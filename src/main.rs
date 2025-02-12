use axum::{routing::get, Router};
use secretsquirrel::physical::{delete, read, write};

#[tokio::main]
async fn main() {
    write("test", "Hello, World!safdasfas1").await;
    println!("{:?}", read("test").await);
    delete("test").await;
    // let app = Router::new().route("/", get(|| async { "Hello, World!" }));
    // let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    // axum::serve(listener, app).await.unwrap();
}
