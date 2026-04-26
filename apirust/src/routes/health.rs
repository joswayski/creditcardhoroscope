use std::sync::Arc;

use axum::{Json, extract::State};
use serde_json::json;
use sqlx::query_scalar;

use crate::app_state::AppState;

pub async fn health(State(state): State<Arc<AppState>>) -> Json<serde_json::Value> {
    let db_status = match query_scalar!("SELECT 1").fetch_one(&state.db).await {
        Ok(_) => "ok".to_string(),
        Err(e) => e.to_string(),
    };

    Json(json!({
        "message":  "Saul Goodman!",
        "docs_url": "https://github.com/joswayski/creditcardhoroscope/blob/main/api/src/main.rs",
        "db": db_status,
    }))
}
