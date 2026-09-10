use std::sync::Arc;
use std::time::Duration;

use axum::{
    extract::{
        Query, State,
        ws::{Message, WebSocket, WebSocketUpgrade},
    },
    http::{HeaderMap, StatusCode},
    response::{IntoResponse, Response},
};
use futures::{SinkExt, StreamExt};
use redis::AsyncCommands;
use serde::Deserialize;
use tracing::{error, info};
use uuid::Uuid;

use crate::dto::{ConversationSignalRequest, ConversationSignalResponse, PubSubSignal};
use crate::state::AppState;

#[derive(Deserialize)]
pub struct JoinQuery {
    id: String,
}

/// Axum handler that upgrades HTTP to WebSocket for call signaling.
/// Extracts memberId from X-User-Id header and conversationId from ?id= query param.
pub async fn join_conversation(
    State(state): State<Arc<AppState>>,
    headers: HeaderMap,
    Query(query): Query<JoinQuery>,
    ws: WebSocketUpgrade,
) -> Response {
    info!("try to make connection");

    let member_id_raw = match headers.get("x-user-id").and_then(|v| v.to_str().ok()) {
        Some(v) => v.to_string(),
        None => {
            error!("missing X-User-Id header");
            return (StatusCode::BAD_REQUEST, "fail to parse").into_response();
        }
    };

    let member_id = match Uuid::parse_str(&member_id_raw) {
        Ok(id) => id,
        Err(e) => {
            error!(
                err = %e,
                member_id_raw = %member_id_raw,
                "fail to parse member id from raw string"
            );
            return (StatusCode::BAD_REQUEST, "fail to parse").into_response();
        }
    };

    let conversation_id = query.id;

    ws.on_upgrade(move |socket| handle_socket(socket, state, member_id, conversation_id))
        .into_response()
}

async fn handle_socket(
    socket: WebSocket,
    state: Arc<AppState>,
    member_id: Uuid,
    conversation_id: String,
) {
    let (ws_sink, mut ws_stream) = socket.split();
    let (tx, rx) = tokio::sync::mpsc::unbounded_channel::<Message>();

    // Spawn a dedicated writer task that owns the WebSocket sink.
    // All writes go through the mpsc channel to avoid lock contention.
    let writer_handle = tokio::spawn({
        let mut sink = ws_sink;
        let mut rx = rx;
        async move {
            while let Some(msg) = rx.recv().await {
                if sink.send(msg).await.is_err() {
                    break;
                }
            }
            let _ = sink.close().await;
        }
    });

    // Register this connection in the local map
    state.conns.insert(member_id, tx.clone());
    info!(
        member_id = %member_id,
        connections = state.conns.len(),
        "success to make connection"
    );

    // Get a multiplexed Valkey connection for commands (SADD, SMEMBERS, SREM, PUBLISH)
    let mut cmd_conn = match state.valkey.get_multiplexed_async_connection().await {
        Ok(c) => c,
        Err(e) => {
            error!(err = %e, "fail to get valkey connection");
            let _ = tx.send(Message::Text(e.to_string().into()));
            close_connection(&state, member_id, tx, writer_handle).await;
            return;
        }
    };

    // --- Handshake: get existing participants ---

    // SMEMBERS {conversationId} → list of participant UUID bytes
    let participant_bytes: Vec<Vec<u8>> = cmd_conn
        .smembers(&conversation_id)
        .await
        .unwrap_or_default();
    let participants: Vec<Uuid> = participant_bytes
        .iter()
        .filter_map(|b| Uuid::from_slice(b).ok())
        .filter(|id| *id != member_id)
        .collect();

    // Send existing participant IDs to the newly connected user
    if !participants.is_empty() {
        let resp = ConversationSignalResponse {
            from_ids: participants.clone(),
            signal: None,
        };
        if let Ok(payload) = serde_json::to_string(&resp)
            && tx.send(Message::Text(payload.into())).is_err()
        {
            error!(member_id = %member_id, "fail to write payload");
            close_connection(&state, member_id, tx, writer_handle).await;
            return;
        }
    }

    // SADD {conversationId} {memberId} — register as participant
    if let Err(e) = cmd_conn
        .sadd::<_, _, ()>(&conversation_id, member_id.as_bytes().as_slice())
        .await
    {
        error!(err = %e, "fail to add participant");
        let _ = tx.send(Message::Text(e.to_string().into()));
        close_connection(&state, member_id, tx, writer_handle).await;
        return;
    }

    // --- Notify existing participants about the new joiner ---
    let channel = format!("signal:{}", conversation_id);
    notify_participants(
        &state,
        &mut cmd_conn,
        &channel,
        member_id,
        &participants,
        None,
    )
    .await;

    // --- Spawn Valkey pub/sub listener for cross-pod signals ---
    let pubsub_handle = {
        let state = state.clone();
        let tx = tx.clone();
        let channel = channel.clone();
        tokio::spawn(async move {
            run_pubsub_listener(state, tx, member_id, &channel).await;
        })
    };

    // --- Spawn ping task (30s interval, mirrors Go's ticker) ---
    let ping_handle = {
        let tx = tx.clone();
        tokio::spawn(async move {
            let mut interval = tokio::time::interval(Duration::from_secs(30));
            loop {
                interval.tick().await;
                if tx.send(Message::Ping(Vec::new().into())).is_err() {
                    error!(member_id = %member_id, "ping failed, client unresponsive");
                    break;
                }
            }
        })
    };

    // --- Main read loop: process WebSocket messages from client ---
    loop {
        match ws_stream.next().await {
            Some(Ok(Message::Text(text))) => {
                let req: ConversationSignalRequest = match serde_json::from_str(&text) {
                    Ok(r) => r,
                    Err(e) => {
                        error!(err = %e, "fail to unmarshalling data");
                        let _ = tx.send(Message::Text("incorrect data".into()));
                        break;
                    }
                };

                notify_participants(
                    &state,
                    &mut cmd_conn,
                    &channel,
                    member_id,
                    &req.to_ids,
                    Some(req.signal),
                )
                .await;
            }
            Some(Ok(Message::Close(_))) => {
                info!(member_id = %member_id, "Connection closed smoothly");
                break;
            }
            Some(Err(e)) => {
                // Check if this is a normal close or an actual error
                let err_str = e.to_string();
                if err_str.contains("close") {
                    info!(member_id = %member_id, err = %e, "Connection closed smoothly");
                } else {
                    error!(member_id = %member_id, err = %e, "read error");
                }
                break;
            }
            None => {
                info!(member_id = %member_id, "Connection closed");
                break;
            }
            _ => {} // Ping/Pong handled automatically by tungstenite
        }
    }

    // --- Cleanup ---
    ping_handle.abort();
    pubsub_handle.abort();

    // SREM {conversationId} {memberId} — remove from participants
    let _: Result<(), redis::RedisError> = cmd_conn
        .srem::<_, _, ()>(&conversation_id, member_id.as_bytes().as_slice())
        .await;

    close_connection(&state, member_id, tx, writer_handle).await;
}

/// Deliver a signal to target participants.
/// Tries local delivery first via the connection map. For participants not found
/// locally (on another pod), publishes to the Valkey pub/sub channel so the
/// subscribing pod can deliver.
async fn notify_participants(
    state: &Arc<AppState>,
    cmd_conn: &mut redis::aio::MultiplexedConnection,
    channel: &str,
    from_id: Uuid,
    to_ids: &[Uuid],
    signal: Option<serde_json::Value>,
) {
    let resp = ConversationSignalResponse {
        from_ids: vec![from_id],
        signal: signal.clone(),
    };
    let resp_raw = match serde_json::to_string(&resp) {
        Ok(r) => r,
        Err(e) => {
            error!(err = %e, "fail to serialize signal response");
            return;
        }
    };

    let mut non_local_ids: Vec<Uuid> = Vec::new();

    for to_id in to_ids {
        if let Some(peer_tx) = state.conns.get(to_id) {
            if peer_tx
                .send(Message::Text(resp_raw.clone().into()))
                .is_err()
            {
                error!(to_id = %to_id, "fail to write payload");
                // Peer's writer is dead; treat as non-local for pub/sub delivery.
                // The peer's own cleanup will remove its map entry.
                non_local_ids.push(*to_id);
            }
        } else {
            non_local_ids.push(*to_id);
        }
    }

    // Publish to Valkey for participants on other pods
    if !non_local_ids.is_empty() {
        let pubsub_msg = PubSubSignal {
            from_id,
            to_ids: non_local_ids,
            signal,
        };
        if let Ok(msg) = serde_json::to_string(&pubsub_msg) {
            let result: Result<(), redis::RedisError> = cmd_conn.publish(channel, &msg).await;
            if let Err(e) = result {
                error!(err = %e, "fail to publish signal");
            }
        }
    }
}

/// Listens to Valkey pub/sub messages on the conversation channel and delivers
/// cross-pod signals to the local WebSocket connection.
async fn run_pubsub_listener(
    state: Arc<AppState>,
    tx: tokio::sync::mpsc::UnboundedSender<Message>,
    member_id: Uuid,
    channel: &str,
) {
    let mut pubsub_conn = match state.valkey.get_async_pubsub().await {
        Ok(c) => c,
        Err(e) => {
            error!(err = %e, "fail to get pubsub connection");
            return;
        }
    };
    if let Err(e) = pubsub_conn.subscribe(channel).await {
        error!(channel = %channel, err = %e, "fail to subscribe");
        return;
    }
    info!(channel = %channel, "subscribed to valkey channel");

    let mut stream = pubsub_conn.on_message();
    while let Some(msg) = stream.next().await {
        let payload: String = match msg.get_payload() {
            Ok(p) => p,
            Err(e) => {
                error!(err = %e, "fail to get pubsub payload");
                continue;
            }
        };
        let signal: PubSubSignal = match serde_json::from_str(&payload) {
            Ok(s) => s,
            Err(e) => {
                error!(err = %e, "fail to parse pubsub signal");
                continue;
            }
        };

        // Skip messages we published ourselves
        if signal.from_id == member_id {
            continue;
        }

        // Only deliver if this connection is an intended recipient
        if !signal.to_ids.contains(&member_id) {
            continue;
        }

        let resp = ConversationSignalResponse {
            from_ids: vec![signal.from_id],
            signal: signal.signal,
        };
        if let Ok(resp_raw) = serde_json::to_string(&resp)
            && tx.send(Message::Text(resp_raw.into())).is_err()
        {
            // Channel closed — the main WS connection is gone
            break;
        }
    }
}

/// Gracefully close a WebSocket connection and remove from the connection map.
async fn close_connection(
    state: &Arc<AppState>,
    member_id: Uuid,
    tx: tokio::sync::mpsc::UnboundedSender<Message>,
    writer_handle: tokio::task::JoinHandle<()>,
) {
    state.conns.remove(&member_id);
    let _ = tx.send(Message::Close(None));
    drop(tx);
    let _ = writer_handle.await;
    info!(
        member_id = %member_id,
        connections = state.conns.len(),
        "success to close connection"
    );
}
