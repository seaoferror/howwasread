use axum::extract::ws::Message;
use dashmap::DashMap;
use uuid::Uuid;

/// Per-connection sender channel.
/// Messages sent through this channel are forwarded to the WebSocket sink
/// by a dedicated writer task.
pub type WsSender = tokio::sync::mpsc::UnboundedSender<Message>;

/// Shared application state across all WebSocket connections.
pub struct AppState {
    /// Local WebSocket connections: memberId → sender channel
    pub conns: DashMap<Uuid, WsSender>,
    /// Valkey client for commands (SADD/SMEMBERS/SREM) and pub/sub
    pub valkey: redis::Client,
}
