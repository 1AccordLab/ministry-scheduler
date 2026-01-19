use axum::{Router, routing::get};

#[actix::main]
async fn main() {
    let addr = "0.0.0.0:3000";
    let app = Router::new().route("/", get(index));
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();

    println!("Server is listening on {}...", addr);
    axum::serve(listener, app).await.unwrap();
}

async fn index() -> String {
    String::from("Hello, world!")
}
