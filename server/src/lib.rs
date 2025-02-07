pub mod oauth2;

use axum::Router;
use dioxus::prelude::{DioxusRouterExt, Element, ServeConfig};
use oauth2::apis::SessionStore;

pub async fn launch(app: fn() -> Element) {
    dotenv::dotenv().ok();

    let addr = dioxus_cli_config::fullstack_address_or_localhost();
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    let app = Router::new()
        .serve_dioxus_application(ServeConfig::new().unwrap(), app)
        .nest("/", oauth2::apis::router())
        .with_state(SessionStore::default());

    axum::serve(listener, app).await.unwrap();
}
