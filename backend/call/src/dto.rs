use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// WebSocket request from client: route a signal to specific users
#[derive(Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ConversationSignalRequest {
    pub to_ids: Vec<Uuid>,
    pub signal: serde_json::Value,
}

/// WebSocket response to client: a signal from specific users
#[derive(Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ConversationSignalResponse {
    pub from_ids: Vec<Uuid>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub signal: Option<serde_json::Value>,
}

/// Internal pub/sub payload for cross-pod signal delivery via Valkey
#[derive(Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct PubSubSignal {
    pub from_id: Uuid,
    pub to_ids: Vec<Uuid>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub signal: Option<serde_json::Value>,
}
