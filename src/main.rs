use axum::{routing::get, Router};
use libsql::Builder;
use libsql::Connection;
use std::time::{SystemTime, UNIX_EPOCH};

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

// CREATE TABLE secrets_a3 (
// id_a3 INTEGER PRIMARY KEY AUTOINCREMENT,
// key_a3 TEXT NOT NULL,
// value_a3 TEXT NOT NULL,
// version_a3 INTEGER DEFAULT (1) NOT NULL,
// updated_at_a3 INTEGER DEFAULT (-1) NOT NULL,
// is_deleted_a3 INTEGER DEFAULT (0) NOT NULL
// );
// CREATE UNIQUE INDEX secrets_a3_key_a3_IDX ON secrets_a3 (key_a3,version_a3);

async fn establish_connection() -> libsql::Result<Connection> {
    let url = std::env::var("TURSO_DATABASE_URL").expect("TURSO_DATABASE_URL must be set");
    let token = std::env::var("TURSO_AUTH_TOKEN").expect("TURSO_AUTH_TOKEN must be set");

    let db = Builder::new_remote(url, token).build().await?;
    db.connect()
}

async fn get_current_version(conn: &Connection, key: &str) -> i64 {
    let mut rows = conn
        .query(
            "SELECT version_a3 FROM secrets_a3 WHERE key_a3 = ? ORDER BY version_a3 DESC LIMIT 1;",
            libsql::params![key],
        )
        .await
        .unwrap();
    if let Some(row) = rows.next().await.unwrap() {
        row.get(0).unwrap()
    } else {
        0
    }
}

async fn read(conn: &Connection, key: &str) -> Option<String> {
    let mut rows = conn
        .query(
            "SELECT value_a3 FROM secrets_a3 WHERE key_a3 = ? AND is_deleted_a3 = 0 ORDER BY version_a3 DESC LIMIT 1;",
            libsql::params![key],
        )
        .await
        .unwrap();
    if let Some(row) = rows.next().await.unwrap() {
        Some(row.get(0).unwrap())
    } else {
        None
    }
}

async fn write(conn: &Connection, key: &str, value: &str) -> () {
    let next_version = get_current_version(conn, key).await + 1;
    let current_epoch_time: i64 = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_secs() as i64;
    conn.execute(
        "INSERT INTO secrets_a3 (key_a3, value_a3, version_a3, updated_at_a3) VALUES (?, ?, ?, ?);",
        libsql::params![key, value, next_version, current_epoch_time],
    )
    .await
    .unwrap();
}

async fn delete(conn: &Connection, key: &str) -> () {
    let current_epoch_time: i64 = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_secs() as i64;
    conn.execute(
        "UPDATE secrets_a3 SET is_deleted_a3 = 1, updated_at_a3 = ? WHERE key_a3 = ?;",
        libsql::params![current_epoch_time, key],
    )
    .await
    .unwrap();
}
