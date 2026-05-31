-- name: CreateCourse :one
INSERT INTO courses (id, code, title, description, instructor_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCourseById :one
SELECT * FROM courses WHERE id = $1;

-- name: ListCourses :many
SELECT * FROM courses
WHERE (sqlc.narg('instructor_id')::text IS NULL OR instructor_id = sqlc.narg('instructor_id'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
