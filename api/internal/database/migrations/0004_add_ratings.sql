ALTER TABLE generations
ADD COLUMN IF NOT EXISTS rating TEXT,
ADD COLUMN IF NOT EXISTS external_id TEXT;

-- is_public BOOLEAN default false,; future PR TODO