use axum::{extract::State, response::Html, routing::get, Router};
use minijinja::{path_loader, Environment};
use serde::Serialize;
use tokio::net::TcpListener;
use tower_http::services::ServeDir;

#[derive(Clone)]
struct AppState {
    minijinja: Environment<'static>,
}

#[tokio::main]
async fn main() {
    let mut env = Environment::new();
    env.set_loader(path_loader("templates"));

    let listener = TcpListener::bind("127.0.0.1:3000").await.unwrap();
    let app = Router::new()
        .route("/", get(index))
        .nest_service("/assets", ServeDir::new("assets"))
        .with_state(AppState { minijinja: env });

    axum::serve(listener, app).await.unwrap();
}

async fn index(State(state): State<AppState>) -> Html<String> {
    let html = state
        .minijinja
        .get_template("index.html")
        .unwrap()
        .render(Person {
            first_name: "Yu Chen".to_string(),
            last_name: "Chung".to_string(),
        })
        .unwrap();

    Html(html)
}

#[derive(Serialize)]
struct Person {
    first_name: String,
    last_name: String,
}
