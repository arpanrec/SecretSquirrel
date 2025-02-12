use libsql::Builder;
use libsql::Connection;
use std::time::{SystemTime, UNIX_EPOCH};

// CREATE TABLE secrets_d (
//     id_d INTEGER PRIMARY KEY AUTOINCREMENT,
//     key_d TEXT NOT NULL,
//     value_d TEXT NOT NULL,
//     version_d INTEGER DEFAULT (1) NOT NULL,
//     updated_at_d INTEGER DEFAULT (-1) NOT NULL,
//     is_deleted_d INTEGER DEFAULT (0) NOT NULL
// );
// CREATE UNIQUE INDEX secrets_d_key_d_IDX ON secrets_d (key_d,version_d);

const TABLE_NAME: &str = "secrets_d";

async fn establish_connection() -> Connection {
    let url = std::env::var("TURSO_DATABASE_URL").expect("TURSO_DATABASE_URL must be set");
    let token = std::env::var("TURSO_AUTH_TOKEN").expect("TURSO_AUTH_TOKEN must be set");

    let db = Builder::new_remote(url, token).build().await.unwrap();
    db.connect().unwrap()
}

async fn get_current_version(conn: &Connection, key: &str) -> i64 {
    let mut rows = conn
        .query(
            &format!(
                "SELECT version_d FROM {TABLE_NAME} WHERE key_d = ? ORDER BY version_d DESC LIMIT 1;"
            ),
            libsql::params![key],).await.unwrap();
    if let Some(row) = rows.next().await.unwrap() {
        row.get(0).unwrap()
    } else {
        0
    }
}

pub async fn read(key: &str) -> Option<String> {
    let mut rows = establish_connection().await
        .query(
            &format!(
                "SELECT value_d FROM {TABLE_NAME} WHERE key_d = ? AND is_deleted_d = 0 ORDER BY version_d DESC LIMIT 1;"
            ),
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

pub async fn write(key: &str, value: &str) -> () {
    let conn = establish_connection().await;
    let next_version = get_current_version(&conn, key).await + 1;
    let current_epoch_time: i64 = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_secs() as i64;
    conn.execute(
        &format!(
            "INSERT INTO {TABLE_NAME} (key_d, value_d, version_d, updated_at_d) VALUES (?, ?, ?, ?);"
        ),
        libsql::params![key, value, next_version, current_epoch_time],
    )
        .await
        .unwrap();
}

pub async fn delete(key: &str) -> () {
    let conn = establish_connection().await;
    let current_epoch_time: i64 = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_secs() as i64;
    conn.execute(
        &format!("UPDATE {TABLE_NAME} SET is_deleted_d = 1, updated_at_d = ? WHERE key_d = ?;"),
        libsql::params![current_epoch_time, key],
    )
    .await
    .unwrap();
}
