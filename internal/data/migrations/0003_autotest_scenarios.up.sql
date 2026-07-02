CREATE TABLE IF NOT EXISTS autotest_scenarios (
    id BIGSERIAL PRIMARY KEY,
    label TEXT NOT NULL UNIQUE,
    user_id TEXT NOT NULL,
    amount INTEGER NOT NULL CHECK (amount > 0),
    ttl_days INTEGER NOT NULL CHECK (ttl_days > 0),
    parallel_requests INTEGER NOT NULL CHECK (parallel_requests >= 2),
    ledger_label TEXT NOT NULL DEFAULT 'test',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
