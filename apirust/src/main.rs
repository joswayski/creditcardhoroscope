use std::{os::unix::net::SocketAddr, time::Duration};

use axum::{
    Router,
    extract::DefaultBodyLimit,
    http::{HeaderValue, Method, StatusCode},
    middleware as axum_middleware,
    routing::get,
};
use tokio::{net::TcpListener, time::Timeout};
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
use utils::rate_limiter;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let config = Config::new()?;

    let cors = CorsLayer::new()
        .allow_methods([Method::GET, Method::POST, Method::PATCH, Method::OPTIONS])
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

    let server = Router::new().merge(general_routes).layer(middleware_layers);

    let address = format!("0.0.0.0:{}", config.api.port);
    let listener = TcpListener::bind(address).await?;

    axum::serve(
        listener,
        server.into_make_service_with_connect_info::<std::net::SocketAddr>(),
    )
    .await?;

    Ok(())
}
