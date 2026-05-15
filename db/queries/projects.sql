-- name: CreateProject :one
INSERT INTO projects (
  id, name, description, student_email, student_name, tech_stack, categories, stars, repo_url
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: GetProject :one
SELECT * FROM projects WHERE id = ?;

-- name: GetProjectByRepoURL :one
SELECT * FROM projects WHERE repo_url = ?;

-- name: GetProjects :many
SELECT * FROM projects ORDER BY created_at DESC;

-- name: GetProjectsByStudent :many
SELECT * FROM projects WHERE student_email = ? ORDER BY created_at DESC;

-- name: UpdateProject :one
UPDATE projects
SET name = ?, description = ?, tech_stack = ?, categories = ?, stars = ?, repo_url = ?
WHERE id = ?
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = ?;

-- name: StarProject :one
UPDATE projects SET stars = stars + 1 WHERE id = ? RETURNING *;
