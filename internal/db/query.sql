-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetUserWithUsername :many
SELECT username, avatar FROM users
WHERE username LIKE $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  user_id, username, email, password, google_id, avatar
) VALUES (
  $1, $2, $3, $4, $5, $6
) ON CONFLICT (username) DO NOTHING
RETURNING *;

-- name: UpdateUser :one
UPDATE users
  SET username = $2,
  email = $3,
  password = $4,
  avatar = $5,
  bio = $6,
  github = $7,
  linkedin = $8,
  website = $9,
  updated_at = CURRENT_TIMESTAMP
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
WHERE name LIKE $1;

-- name: CreateRole :one
INSERT INTO roles (
  role_id, name
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateRole :one
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

-- name: GetAllPosts :many
SELECT * FROM posts;

-- name: GetPost :one
SELECT * FROM posts
WHERE post_id = $1;

-- name: CreatePost :one
INSERT INTO posts (
  post_id, user_id, title, content, photo, video
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdatePost :one
UPDATE posts
  SET title = $3,
  content = $4,
  photo = $5,
  video = $6,
  updated_at = CURRENT_TIMESTAMP
WHERE post_id = $1 AND
user_id = $2
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE post_id = $1;

-- name: GetAllUserPosts :many
SELECT * FROM posts
WHERE user_id = $1;

-- name: GetUserPost :one
SELECT * FROM posts
WHERE user_id = $1 AND
post_id = $2;

-- name: GetAllComments :many
SELECT * FROM comments
WHERE post_id = $1;

-- name: GetComment :one
SELECT * FROM comments
WHERE comment_id = $1;

-- name: CreateComment :one
INSERT INTO comments (
  comment_id, user_id, post_id, content
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateComment :one
UPDATE comments
  SET content = $2,
  updated_at = CURRENT_TIMESTAMP
WHERE comment_id = $1 AND
user_id = $3
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE comment_id = $1 AND
user_id = $2
;

-- name: CreateLike :one
INSERT INTO likes (
  user_id, post_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteLike :exec
DELETE FROM likes
WHERE user_id = $1 AND
post_id = $2;

-- name: GetLike :one
SELECT * FROM likes
WHERE user_id = $1 AND
post_id = $2;