use std::sync::Arc;

use axum::Router;
use axum::routing::get;
use dashmap::DashMap;
use redis::{Client, TlsCertificates};
use std::{env, fs};
use tracing::info;

mod dto;
mod signal;
mod state;

use state::AppState;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt().json().init();

    let client = Client::build_with_tls(
        format!(
            "rediss://{}:{}@{}",
            env::var("VALKEY_USERNAME").unwrap(),
            env::var("VALKEY_PASSWORD").unwrap(),
            env::var("VALKEY_ADDRESS").unwrap()
        ),
        TlsCertificates {
            client_tls: None,
            root_cert: Some(fs::read(env::var("VALKEY_CA_CERT_PATH").unwrap()).unwrap()),
        },
    )
    .unwrap();

    let state = Arc::new(AppState {
        conns: DashMap::new(),
        valkey: client,
    });

    let app = Router::new()
        .route("/call/join", get(signal::join_conversation))
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    info!("listening on http://0.0.0.0:3000");
    axum::serve(listener, app).await.unwrap();
}
