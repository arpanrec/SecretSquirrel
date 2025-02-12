use axum::http::StatusCode;
use axum::{extract::Path, http::Response, routing::get, Router};
use secretsquirrel::physical::{delete, read, write};

#[tokio::main]
async fn main() {
    let app = Router::new().route(
        "/secret/{*key}",
        get(handle_get).post(handle_post).delete(handle_delete),
    );
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

async fn handle_get(Path(path): Path<String>) -> Response<String> {
    println!("Received request for key: {}", path);
    let value = read(&path).await;
    println!("Read value: {:?}", value);
    let builder = Response::builder();
    if let Some(value) = value {
        println!("Found value: {}", value);
        builder
            .status(StatusCode::OK)
            .header("Content-Type", "text/plain")
            .body(value)
            .expect("Failed to send response")
    } else {
        println!("Not Found");
        builder
            .status(StatusCode::NOT_FOUND)
            .header("Content-Type", "text/plain")
            .body("".to_string())
            .expect("Failed to send response")
    }
}

async fn handle_post(Path(path): Path<String>, body: String) -> Response<String> {
    println!("Received request for key: {}", path);
    println!("Received body: {}", body);
    write(&path, &body).await;
    Response::new("".to_string())
}

async fn handle_delete(Path(path): Path<String>) -> Response<String> {
    println!("Received request for key: {}", path);
    delete(&path).await;
    Response::new("".to_string())
}
