-- name: CreateUser :one
-- Creates a new user and returns the created record.
INSERT INTO users (name, dob)
VALUES ($1, $2)
RETURNING id, name, dob;

-- name: GetUser :one
-- Fetches a single user by their primary key.
SELECT id, name, dob
FROM users
WHERE id = $1;

-- name: UpdateUser :one
-- Updates name and dob for an existing user, returns updated record.
UPDATE users
SET name = $1,
    dob  = $2
WHERE id = $3
RETURNING id, name, dob;

-- name: DeleteUser :exec
-- Permanently removes a user by ID.
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
-- Returns all users ordered by ID, with optional pagination.
SELECT id, name, dob
FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
-- Returns total count of users (used for pagination metadata).
SELECT COUNT(*) FROM users;
