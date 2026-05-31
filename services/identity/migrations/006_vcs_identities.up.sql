CREATE TABLE vcs_identities(
    user_id TEXT NOT NULL,
    provider_kind INT NOT NULL,
    provider_instance TEXT NOT NULL DEFAULT '',
    external_user_id BIGINT NOT NULL,
    external_login TEXT NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP,
    token_valid BOOLEAN NOT NULL DEFAULT TRUE,
    linked_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT vcs_identities_pkey PRIMARY KEY (user_id, provider_kind, provider_instance)
);

CREATE INDEX vcs_identities_user_idx ON vcs_identities (user_id);
