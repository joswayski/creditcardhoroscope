use axum::{Router, routing::get};
use tokio::net::TcpListener;

mod config;
use config::Config;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();
    let config = Config::new()?;

    let server = Router::new().route("/", get(root));

    let address = format!("0.0.0.0:{}", config.api.port);
    let listener = TcpListener::bind(address).await?;

    axum::serve(listener, server).await?;

    Ok(())
}

async fn root() -> String {
    String::from("Hello")
}
