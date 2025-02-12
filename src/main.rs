use axum::{extract::Path, http::Response, routing::get, Router};
use axum::http::StatusCode;
use secretsquirrel::physical::{delete, read, write};

#[tokio::main]
async fn main() {
    // write("test", "Hello, World!safdasfas1").await;
    // println!("{:?}", read("test").await);
    // delete("test").await;
    let app = Router::new().route(
        "/secret/{*key}",
        get(|Path(path): Path<String>| async move {
            // Response::builder()
            //     .status(StatusCode::OK)
            //     .header("Content-Type", "text/plain")
            //     .body(read(&path).await.unwrap())
            //     .unwrap()
            println!("Received request for key: {}", path);
            let value = read(&path).await;
            println!("Read value: {:?}", value);
            if value.is_some() {
                let value = value.unwrap();
                println!("Found value: {}", value);
                Response::builder()
                    .status(StatusCode::OK)
                    .header("Content-Type", "text/plain")
                    .body(value)
                    .expect("Failed to send response");
            } else {
                println!("Not Found");
                Response::builder()
                    .status(404)
                    .header("Content-Type", "text/plain")
                    .body("")
                    .expect("Failed to send response");
            };
        }),
    );
    // let app: Router = Router::new().route("/{*key}", get(handler));
    //
    // async fn handler(Path(path): Path<String>) -> String {
    //     path
    // }
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
