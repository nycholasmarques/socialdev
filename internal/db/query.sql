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
SELECT * FROM posts
WHERE (title ILIKE '%' || sqlc.narg('title') || '%' OR sqlc.narg('title') IS NULL)
AND (content ILIKE '%' || sqlc.narg('content') || '%' OR sqlc.narg('content') IS NULL);

-- name: GetAllPostsWithRelations :many
SELECT 
    p.post_id,
    p.user_id,
    p.title,
    p.content,
    p.photo,
    p.video,
    p.created_at,
    p.updated_at,
    u.username AS user_name,
    u.avatar AS user_avatar,
    COALESCE(l.likes_count, 0) AS like_count,
    COALESCE(v.views_count, 0) AS views_count,
    c.comment_id,
    c.content AS comment_content,
    c.created_at AS comment_created_at,
    cu.user_id AS comment_user_id,
    cu.username AS comment_user_name,
    cu.avatar AS comment_user_avatar
FROM posts p
LEFT JOIN users u ON p.user_id = u.user_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS likes_count FROM likes GROUP BY post_id
) l ON p.post_id = l.post_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS views_count FROM views GROUP BY post_id
) v ON p.post_id = v.post_id
LEFT JOIN comments c ON c.post_id = p.post_id
LEFT JOIN users cu ON c.user_id = cu.user_id
WHERE (p.title ILIKE '%' || sqlc.narg('title') || '%' OR sqlc.narg('title') IS NULL)
AND (p.content ILIKE '%' || sqlc.narg('content') || '%' OR sqlc.narg('content') IS NULL);

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