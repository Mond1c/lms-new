-- name: CreateEnrollment :one
INSERT INTO enrollments (id, user_id, course_id, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteEnrollment :execrows
DELETE FROM enrollments WHERE user_id = $1 AND course_id = $2;

-- name: ListEnrollments :many
SELECT * FROM enrollments
WHERE (sqlc.narg('course_id')::text IS NULL OR course_id = sqlc.narg('course_id'))
  AND (sqlc.narg('user_id')::text IS NULL OR user_id = sqlc.narg('user_id'))
ORDER BY enrolled_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
