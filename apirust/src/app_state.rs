use std::sync::Arc;

use sqlx::{PgPool, Pool, Postgres, postgres::PgPoolOptions};

use crate::{
    app_state,
    config::{Config, DatabaseConfig},
};

pub struct AppState {
    db: PgPool,
}

impl AppState {
    pub async fn new(config: &Config) -> Result<Arc<AppState>, Vec<String>> {
        let db_pool = get_db(&config.database).await;

        match db_pool {
            Ok(db_pool) => Ok(Arc::new(AppState { db: db_pool })),
            (db_pool) => {
                let errors = [db_pool.err()].into_iter().flatten().collect();
                Err(errors)
            }
        }
    }
}

async fn get_db(dbConfig: &DatabaseConfig) -> Result<PgPool, String> {
    PgPoolOptions::new()
        .max_connections(10)
        .min_connections(2)
        .connect(dbConfig.database_url.as_str())
        .await
        .map_err(|e| format!("Error ocurred connecting to database: {e}"))
}
