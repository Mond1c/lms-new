-- name: RegisterStudentRepo :one
INSERT INTO student_repos (
    id, user_id, assignment_id, provider_kind, provider_instance,
    full_name, external_id, state, clone_url_https, clone_url_ssh
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
ON CONFLICT (user_id, assignment_id) DO UPDATE SET
    provider_kind = EXCLUDED.provider_kind,
    provider_instance = EXCLUDED.provider_instance,
    full_name = EXCLUDED.full_name,
    external_id = EXCLUDED.external_id,
    state = EXCLUDED.state,
    clone_url_https = EXCLUDED.clone_url_https,
    clone_url_ssh = EXCLUDED.clone_url_ssh
RETURNING *;

-- name: GetStudentRepo :one
SELECT * FROM student_repos WHERE user_id = $1 AND assignment_id = $2;
