use std::sync::Arc;

use sqlx::{PgPool, postgres::PgPoolOptions};

use crate::config::{Config, DatabaseConfig};

pub struct AppState {
    pub db: PgPool,
}

impl AppState {
    pub async fn new(config: &Config) -> anyhow::Result<Arc<AppState>> {
        let db_pool = get_db(&config.database).await;

        match db_pool {
            Ok(db_pool) => Ok(Arc::new(AppState { db: db_pool })),
            (db_pool) => {
                let errors = [db_pool.err()]
                    .into_iter()
                    .flatten()
                    .collect::<Vec<String>>();

                Err(anyhow::anyhow!("App state error:\n{}", errors.join("\n")))
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
