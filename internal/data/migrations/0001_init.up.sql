CREATE TABLE IF NOT EXISTS points_lots (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    amount INTEGER NOT NULL CHECK (amount > 0),
    remaining INTEGER NOT NULL CHECK (remaining >= 0),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_points_lots_user_expiry ON points_lots (user_id, expires_at);

CREATE TABLE IF NOT EXISTS holds (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    amount INTEGER NOT NULL CHECK (amount > 0),
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','confirmed','cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_holds_user ON holds (user_id, status);

CREATE TABLE IF NOT EXISTS hold_allocations (
    id BIGSERIAL PRIMARY KEY,
    hold_id BIGINT NOT NULL REFERENCES holds (id),
    lot_id BIGINT NOT NULL REFERENCES points_lots (id),
    amount INTEGER NOT NULL CHECK (amount > 0)
);
CREATE INDEX IF NOT EXISTS idx_hold_allocations_hold ON hold_allocations (hold_id);

CREATE TABLE IF NOT EXISTS ledger_entries (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    amount INTEGER NOT NULL,
    ref_type TEXT,
    ref_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_ledger_user ON ledger_entries (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS idempotency_keys (
    key TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    response_status INT NOT NULL DEFAULT 0,
    response_body JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (key, endpoint)
);
