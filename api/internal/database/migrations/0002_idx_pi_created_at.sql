CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_generations_pi_created 
ON generations(payment_intent_id, created_at DESC);