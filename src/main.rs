use axum::{routing::get, Router};
use secretsquirrel::physical::{delete, establish_connection, read, write};

#[tokio::main]
async fn main() {
    let conn = establish_connection()
        .await
        .expect("Failed to connect to database");
    write(&conn, "test", "Hello, World!safdasfas1").await;
    println!("{:?}", read(&conn, "test").await);
    delete(&conn, "test").await;
    // let app = Router::new().route("/", get(|| async { "Hello, World!" }));
    // let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    // axum::serve(listener, app).await.unwrap();
}
