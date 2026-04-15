CREATE TABLE IF NOT EXISTS payment_intents (
    -- Metadata
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Stripe Data
    payment_intent_id TEXT NOT NULL, 
    amount INT NOT NULL,                      -- Amount charged in cents
    currency TEXT NOT NULL DEFAULT 'usd',
    status TEXT NOT NULL DEFAULT 'pending',   -- 'pending', 'processing', 'paid', 'refunded'
    
    -- Card Details (from Stripe)
    card_brand TEXT,
    card_exp_month TEXT,
    card_exp_year TEXT,
    card_last_4 TEXT,
    card_country TEXT,
    card_postal TEXT
);

-- Table 2: Generations (many per payment intent)
CREATE TABLE IF NOT EXISTS generations (
    -- Metadata
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- FK to payment
    payment_intent_id BIGINT NOT NULL REFERENCES payment_intents(id),
    
    -- Generation status and errors
    status TEXT NOT NULL DEFAULT 'pending',  -- 'pending', 'completed', 'failed'
    error TEXT,
    
    -- AI/OpenRouter specific
    or_gen_id TEXT,
    or_model TEXT,
    or_tokens_used INT,
    horoscope TEXT
);