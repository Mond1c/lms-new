-- name: UpsertVCSIdentity :one
INSERT INTO vcs_identities (
    user_id, provider_kind, provider_instance, external_user_id, external_login,
    access_token, refresh_token, expires_at, token_valid
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, TRUE
)
ON CONFLICT (user_id, provider_kind, provider_instance) DO UPDATE SET
    external_user_id = EXCLUDED.external_user_id,
    external_login = EXCLUDED.external_login,
    access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token,
    expires_at = EXCLUDED.expires_at,
    token_valid = TRUE,
    linked_at = now()
RETURNING *;

-- name: DeleteVCSIdentity :execrows
DELETE FROM vcs_identities
WHERE user_id = $1 AND provider_kind = $2 AND provider_instance = $3;

-- name: ListVCSIdentities :many
SELECT * FROM vcs_identities WHERE user_id = $1 ORDER BY linked_at DESC;
