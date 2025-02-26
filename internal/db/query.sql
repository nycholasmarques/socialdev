-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT username, avatar, bio, github, linkedin, website, email, created_at FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetUserWithUsername :one
SELECT username, avatar FROM users
WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (
  user_id, username, email, password
) VALUES (
  $1, $2, $3, $4
) ON CONFLICT (username) DO NOTHING
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
  SET username = $2,
  email = $3,
  password = $4,
  avatar = $5,
  bio = $6,
  github = $7,
  linkedin = $8,
  website = $9
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: GetAllRoles :many
SELECT * FROM roles;

-- name: GetRole :one
SELECT * FROM roles
WHERE role_id = $1;

-- name: GetRoleWithName :many
SELECT * FROM roles
WHERE name = $1;

-- name: CreateRole :one
INSERT INTO roles (
  role_id, name
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateRole :exec
UPDATE roles
  SET name = $2
WHERE role_id = $1
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE role_id = $1;

-- name: CreateUserWithRoles :one
INSERT INTO user_roles (
  user_id, role_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteUserWithRoles :exec
DELETE FROM user_roles 
WHERE user_id = $1 AND
role_id = $2;

