use axum::{Json, response::IntoResponse};
use serde_json::json;

pub async fn root() -> Json<serde_json::Value> {
    Json(json!({
        "message":  "Hi! You probably meant to go to one of the other routes. Make sure to check the documentation!",
        "docs_url": "https://github.com/joswayski/creditcardhoroscope/blob/main/api/src/main.rs"
    }))
}
