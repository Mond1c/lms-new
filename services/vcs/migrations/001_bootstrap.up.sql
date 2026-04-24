-- M1 placeholder for vcs-svc schema. Real tables land later.
CREATE TABLE IF NOT EXISTS schema_bootstrap (
    id          SMALLINT PRIMARY KEY DEFAULT 1,
    initialized BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT schema_bootstrap_singleton CHECK (id = 1)
);
INSERT INTO schema_bootstrap (id) VALUES (1) ON CONFLICT (id) DO NOTHING;
