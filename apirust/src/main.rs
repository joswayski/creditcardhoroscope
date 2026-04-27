use std::{os::unix::net::SocketAddr, sync::Arc, time::Duration};

use axum::{
    Router,
    extract::DefaultBodyLimit,
    http::{HeaderValue, Method, StatusCode},
    middleware as axum_middleware,
    routing::get,
};
use sqlx::{PgPool, postgres::PgPoolOptions};
use tokio::{net::TcpListener, signal, time::Timeout};
use tower::{
    ServiceBuilder,
    limit::{RateLimitLayer, rate},
};
use tower_governor::{
    GovernorLayer,
    governor::{GovernorConfig, GovernorConfigBuilder},
    key_extractor::SmartIpKeyExtractor,
};
use tower_http::trace::TraceLayer;
use tower_http::{
    cors::{Any, CorsLayer},
    timeout::TimeoutLayer,
};

mod config;
use config::Config;

mod routes;
use routes::{health, root};

mod utils;
use tracing::{error, info};
use utils::rate_limiter;

mod database;
use database::run_migrations;

use crate::{app_state::AppState, config::AppConfig};

mod app_state;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let config = Config::new()?;

    let app_state = AppState::new(&config).await?;

    let cors = CorsLayer::new()
        .allow_methods(config.server.allowed_methods)
        .allow_origin(config.server.allowed_origins);

    let middleware_layers = ServiceBuilder::new()
        .layer(TraceLayer::new_for_http())
        .layer(TimeoutLayer::with_status_code(
            // second
            StatusCode::REQUEST_TIMEOUT,
            Duration::from_secs(config.server.timeout_seconds),
        ))
        .layer(cors)
        .layer(DefaultBodyLimit::max(config.server.max_body_bytes));

    let general_routes = Router::new()
        .route("/", get(root))
        .route("/api/v1", get(root))
        .route("/api/v1/health", get(health))
        .fallback(root)
        .layer(GovernorLayer::new(rate_limiter::new(1, 10))); // General, applies to all

    let server = Router::new()
        .merge(general_routes)
        .layer(middleware_layers)
        .with_state(app_state);

    let address = format!("0.0.0.0:{}", config.server.port);
    let listener = TcpListener::bind(address).await?;

    if let Err(e) = axum::serve(
        listener,
        server.into_make_service_with_connect_info::<std::net::SocketAddr>(),
    )
    .with_graceful_shutdown(shutdown_signal())
    .await
    {
        error!("Server error: {}", e);
        return Err(e.into());
    };

    Ok(())
}

async fn shutdown_signal() {
    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install ctrl_c handler")
    };

    #[cfg(unix)]
    let terminate = async {
        use tokio::signal::unix::SignalKind;

        signal::unix::signal(SignalKind::terminate())
            .expect("failed to install signal handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending();

    tokio::select! {
        () = ctrl_c => {},
        () = terminate => {},
    }

    info!("shutdown signal received");
}
