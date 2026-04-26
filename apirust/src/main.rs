use axum::{Router, routing::get};
use tokio::net::TcpListener;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::fmt::init();

    let server = Router::new().route("/", get(root));

    let listener = TcpListener::bind("0.0.0.0:8080").await?;

    axum::serve(listener, server).await?;

    Ok(())
}

async fn root() -> String {
    String::from("Hello")
}
