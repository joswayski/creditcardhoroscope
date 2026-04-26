use axum::http::HeaderValue;
use tracing::warn;

const ALLOWED_ORIGINS: [&str; 3] = [
    "https://creditcardhoroscope.com",
    "https://staging.creditcardhoroscope.com",
    "http://localhost:5173",
];

pub struct Config {
    pub api: APIConfig,
    pub server: ServerConfig,
    pub stripe: StripeConfig,
    pub ai: AIConfig,
}

pub struct APIConfig {
    pub port: String,
    base_url: url::Url,
    environment: String,
    database_url: url::Url,
    support_email: String,
    max_horoscope_limit: u8,
}

pub struct ServerConfig {
    pub allowed_origins: [HeaderValue; 3],
    pub max_body_bytes: usize,
    pub timeout_seconds: u64,
}

struct AIConfig {
    base_url: url::Url,
    api_key: String,
    model: String,
    system_prompt: String,
}

struct StripeConfig {
    secret_key: String,
    webhook_secret_key: String,
}

impl Config {
    pub fn new() -> anyhow::Result<Config> {
        if let Err(e) = dotenvy::dotenv() {
            tracing::warn!("error" = ?e, "Could not load .env! Continuing...")
        }

        // Get all our configs
        let api_result = APIConfig::new();
        let stripe_result = StripeConfig::new();
        let ai_result = AIConfig::new();
        let server_result = ServerConfig::new();

        match (api_result, stripe_result, ai_result, server_result) {
            (Ok(api), Ok(stripe), Ok(ai), Ok(server)) => Ok(Config {
                api,
                ai,
                stripe,
                server,
            }),
            (api_res, stripe_res, ai_res, server_res) => {
                let config_errors = [
                    api_res.err(),
                    stripe_res.err(),
                    ai_res.err(),
                    server_res.err(),
                ]
                .into_iter()
                .flatten()
                .flatten()
                .collect::<Vec<String>>()
                .join("\n");

                Err(anyhow::anyhow!("Config error: \n{}", config_errors))
            }
        }
    }
}

impl APIConfig {
    fn new() -> Result<APIConfig, Vec<String>> {
        let api_port = get_non_empty("API_PORT").unwrap_or("8080".into());
        let base_url = get_non_empty("BASE_URL")
            .ok_or_else(|| "BASE_URL is required and cannot be empty".into())
            .and_then(|v| {
                url::Url::parse(&v).map_err(|_| {
                    "BASE_URL is invalid. This is used for the shareable links.".into()
                })
            });

        let environment = get_non_empty("ENVIRONMENT").unwrap_or("development".into());
        let database_url = get_non_empty("DATABASE_URL")
            .ok_or_else(|| "DATABASE_URL is required and cannot be empty".into())
            .and_then(|v| {
                url::Url::parse(&v).map_err(|_| {
                    "DATABASE_URL is invalid. You need this to connect to the database.".into()
                })
            });

        let support_email = get_non_empty("SUPPORT_EMAIL")
            .filter(|v| v.contains("@"))
            .unwrap_or("contact@josevalerio.com".into());

        let max_horoscope_limit: u8 = get_non_empty("MAX_HOROSCOPE_LIMIT")
            .and_then(|v| v.parse::<u8>().ok())
            .filter(|&n| n > 0)
            .unwrap_or_else(|| {
                tracing::warn!("MAX_HOROSCOPE_LIMIT is not valid. Defaulting to 3.");
                3
            });

        match (base_url, database_url) {
            (Ok(base_url), Ok(database_url)) => Ok(APIConfig {
                environment: environment,
                port: api_port,
                base_url: base_url,
                database_url: database_url,
                support_email: support_email,
                max_horoscope_limit: max_horoscope_limit,
            }),
            (base_url, database_url) => {
                let errors = [base_url.err(), database_url.err()]
                    .into_iter()
                    .flatten()
                    .collect();
                Err(errors)
            }
        }
    }
}

impl StripeConfig {
    fn new() -> Result<StripeConfig, Vec<String>> {
        let secret_key = get_non_empty("STRIPE_SECRET_KEY")
            .ok_or_else(|| "STRIPE_SECRET_KEY cannot be empty".into())
            .and_then(|v| {
                if v.starts_with("sk_") {
                    Ok(v)
                } else {
                    Err("STRIPE_SECRET_KEY must start with sk_".into())
                }
            });

        let webhook_secret_key = get_non_empty("STRIPE_WEBHOOK_SECRET_KEY")
            .ok_or_else(|| "STRIPE_WEBHOOK_SECRET_KEY cannot be empty".into())
            .and_then(|v| {
                if v.starts_with("whsec_") {
                    Ok(v)
                } else {
                    Err("STRIPE_WEBHOOK_SECRET_KEY must start with whsec_".into())
                }
            });

        match (secret_key, webhook_secret_key) {
            (Ok(secret_key), Ok(webhook_secret_key)) => Ok(StripeConfig {
                secret_key,
                webhook_secret_key,
            }),
            (secret_key, webhook_secret_key) => {
                let errors = [secret_key.err(), webhook_secret_key.err()]
                    .into_iter()
                    .flatten()
                    .collect();
                Err(errors)
            }
        }
    }
}

impl AIConfig {
    fn new() -> Result<AIConfig, Vec<String>> {
        let base_url = get_non_empty("AI_BASE_URL")
            .ok_or_else(|| "AI_BASE_URL cannot be empty".into())
            .and_then(|v| url::Url::parse(&v).map_err(|e| "AI_BASE_URL is an invalid URL".into()));

        let api_key = get_non_empty("AI_API_KEY")
            .ok_or_else(|| "AI_API_KEY cannot be empty".into())
            .and_then(|v| {
                if v.starts_with("sk-or-v1-") {
                    // Only support OR for now
                    Ok(v)
                } else {
                    Err("AI_API_KEY must start with 'sk-or-v1-'".into())
                }
            });

        let model = get_non_empty("AI_MODEL").unwrap_or_else(|| {
            warn!("AI_MODEL not set, defaulting to \"google/gemini-3.1-flash-lite-preview\"");
            "google/gemini-3.1-flash-lite-preview".to_string()
        });

        let system_prompt = get_non_empty("AI_SYSTEM_PROMPT")
            .ok_or_else(|| "AI_SYSTEM_PROMPT cannot be empty".into());

        match (base_url, api_key, model, system_prompt) {
            (Ok(base_url), Ok(api_key), model, Ok(system_prompt)) => Ok(AIConfig {
                base_url,
                api_key,
                model,
                system_prompt: system_prompt,
            }),
            (base_url, api_key, _model, system_prompt) => {
                let errors = [base_url.err(), api_key.err(), system_prompt.err()]
                    .into_iter()
                    .flatten()
                    .collect();

                Err(errors)
            }
        }
    }
}
fn get_non_empty(key: &'static str) -> Option<String> {
    std::env::var(key).ok().filter(|v| !v.is_empty())
}

impl ServerConfig {
    fn new() -> Result<ServerConfig, Vec<String>> {
        let origins = ALLOWED_ORIGINS.map(|v| {
            v.parse::<HeaderValue>()
                .map_err(|e| format!("Could not parse domain to HeaderValue {v} - error: {e}"))
        });

        match origins {
            [Ok(a), Ok(b), Ok(c)] => Ok(ServerConfig {
                allowed_origins: [a, b, c],
                max_body_bytes: 1024,
                timeout_seconds: 30,
            }),
            origins => {
                let errors = origins.into_iter().filter_map(Result::err).collect();

                Err(errors)
            }
        }
    }
}
