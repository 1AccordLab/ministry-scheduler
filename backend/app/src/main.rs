use axum::{response::IntoResponse, routing::get, Json, Router};
use serde::{Deserialize, Serialize};
use tower_http::trace::{DefaultMakeSpan, DefaultOnRequest, DefaultOnResponse, TraceLayer};
use tracing::Level;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt().init();

    let app = Router::new().route("/", get(root)).layer(
        TraceLayer::new_for_http()
            .make_span_with(DefaultMakeSpan::new().level(Level::INFO))
            .on_request(DefaultOnRequest::new().level(Level::INFO))
            .on_response(DefaultOnResponse::new().level(Level::INFO)),
    );

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

#[axum::debug_handler]
async fn root() -> impl IntoResponse {
    Json(HelloWorld {
        msg: "Hello, World!".to_string(),
    })
}

#[derive(Serialize, Deserialize)]
struct HelloWorld {
    msg: String,
}
