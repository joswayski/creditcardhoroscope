use std::time::Duration;

use axum::{
    Router,
    extract::DefaultBodyLimit,
    http::{HeaderValue, Method, StatusCode},
    routing::get,
};
use tokio::{net::TcpListener, time::Timeout};
use tower::ServiceBuilder;
use tower_http::trace::TraceLayer;
use tower_http::{
    cors::{Any, CorsLayer},
    timeout::TimeoutLayer,
};

mod config;
use config::Config;

mod routes;
use routes::{health, root};

const ALLOWED_ORIGINS: [&str; 3] = [
    "https://creditcardhoroscope.com",
    "https://staging.creditcardhoroscope.com",
    "http://localhost:5173",
];

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let config = Config::new()?;

    let cors = CorsLayer::new()
        .allow_methods([Method::GET, Method::POST, Method::PATCH, Method::OPTIONS])
        .allow_origin(ALLOWED_ORIGINS.map(|origin| {
            origin.parse::<HeaderValue>().unwrap_or_else(|_| {
                panic!("allowed origin should be a valid header value: {origin}")
            })
        }));

    let middleware = ServiceBuilder::new()
        .layer(TraceLayer::new_for_http()) // first
        .layer(TimeoutLayer::with_status_code(
            // second
            StatusCode::REQUEST_TIMEOUT,
            Duration::from_secs(30),
        ))
        .layer(cors) // third
        .layer(DefaultBodyLimit::max(1024)); // fourth

    let server = Router::new()
        .layer(middleware)
        .route("/", get(root))
        .route("/api/v1", get(root))
        .route("/api/v1/health", get(health))
        .fallback(root);

    let address = format!("0.0.0.0:{}", config.api.port);
    let listener = TcpListener::bind(address).await?;

    axum::serve(listener, server).await?;

    Ok(())
}
