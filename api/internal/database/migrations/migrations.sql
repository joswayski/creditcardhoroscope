CREATE TABLE IF NOT EXISTS migrations (
    name TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
)