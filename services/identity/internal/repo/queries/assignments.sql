-- name: CreateAssignment :one
INSERT INTO assignments (
    id, course_id, slug, title, description_markdown,
    deadline, hard_deadline, max_score, template_repo, repo_naming_pattern,
    auto_request_review_on_pass, requires_defense,
    weight_tests, weight_quality, defence_multiplier, custom_formula, runner
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10,
    $11, $12,
    $13, $14, $15, $16, $17
)
RETURNING *;

-- name: GetAssignmentById :one
SELECT * FROM assignments WHERE id = $1;

-- name: ListAssignments :many
SELECT * FROM assignments
WHERE (sqlc.narg('course_id')::text IS NULL OR course_id = sqlc.narg('course_id'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
