use governor::middleware::NoOpMiddleware;
use tower_governor::{
    governor::{GovernorConfig, GovernorConfigBuilder},
    key_extractor::SmartIpKeyExtractor,
};

pub fn new(per_second: u64, burst: u32) -> GovernorConfig<SmartIpKeyExtractor, NoOpMiddleware> {
    GovernorConfigBuilder::default()
        .key_extractor(SmartIpKeyExtractor)
        .per_second(per_second)
        .burst_size(burst)
        .finish()
        .expect("Invalid general governor config")
}
