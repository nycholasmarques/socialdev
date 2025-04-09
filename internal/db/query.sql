-- name: GetAllUsers :many
SELECT * FROM users
WHERE (username ILIKE '%' || sqlc.narg('username') || '%' OR sqlc.narg('username') IS NULL)
ORDER BY created_at DESC
LIMIT sqlc.narg('limit')
OFFSET sqlc.narg('offset');

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
    COALESCE(c.comments_count, 0) AS comment_count,
    EXISTS (
        SELECT 1 FROM likes 
        WHERE likes.post_id = p.post_id 
        AND likes.user_id = sqlc.narg('user_id')
    ) AS liked_by_user, -- Adiciona se o usuário curtiu
    COUNT(*) OVER() AS total_count
FROM posts p
LEFT JOIN users u ON p.user_id = u.user_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS likes_count 
    FROM likes 
    GROUP BY post_id
) l ON p.post_id = l.post_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS views_count 
    FROM views 
    GROUP BY post_id
) v ON p.post_id = v.post_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS comments_count 
    FROM comments 
    GROUP BY post_id
) c ON p.post_id = c.post_id
WHERE (p.title ILIKE '%' || sqlc.narg('title') || '%' OR sqlc.narg('title') IS NULL)
AND (p.content ILIKE '%' || sqlc.narg('content') || '%' OR sqlc.narg('content') IS NULL)
ORDER BY 
    CASE 
        WHEN sqlc.narg('order_by') = 'likes' THEN COALESCE(l.likes_count, 0)
        WHEN sqlc.narg('order_by') = 'comments' THEN COALESCE(c.comments_count, 0)
        ELSE 0
    END DESC,
    p.created_at DESC
LIMIT sqlc.narg('limit')
OFFSET sqlc.narg('offset');


-- name: GetPost :one
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
    COALESCE(c.comments_count, 0) AS comment_count,
    EXISTS (
        SELECT 1 FROM likes 
        WHERE likes.post_id = p.post_id 
        AND likes.user_id = sqlc.narg('user_id')
    ) AS liked_by_user, 
    COUNT(*) OVER() AS total_count
FROM posts p
LEFT JOIN users u ON p.user_id = u.user_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS likes_count 
    FROM likes 
    GROUP BY post_id
) l ON p.post_id = l.post_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS views_count 
    FROM views 
    GROUP BY post_id
) v ON p.post_id = v.post_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS comments_count 
    FROM comments 
    GROUP BY post_id
) c ON p.post_id = c.post_id
WHERE p.post_id = $1;

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
WHERE p.user_id = $1
ORDER BY p.created_at DESC;

-- name: GetUserPost :one
SELECT * FROM posts
WHERE user_id = $1 AND
post_id = $2;

-- name: GetAllComments :many
SELECT 
    c.comment_id,
    c.user_id,
    c.post_id,
    c.content,
    c.created_at,
    c.updated_at,
    u.username,
    u.avatar
FROM comments c
JOIN users u ON c.user_id = u.user_id
WHERE c.post_id = $1
ORDER BY c.created_at ASC;

-- name: GetComment :one
SELECT 
    c.comment_id,
    c.user_id,
    c.post_id,
    c.content,
    c.created_at,
    c.updated_at,
    u.username,
    u.avatar
FROM comments c
JOIN users u ON c.user_id = u.user_id
WHERE c.comment_id = $1;

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

-- name: AddFavorite :one
INSERT INTO favorites (user_id, post_id)
VALUES ($1, $2)
RETURNING *;

-- name: RemoveFavorite :exec
DELETE FROM favorites
WHERE user_id = $1 AND post_id = $2;

-- name: GetUserFavorites :many
SELECT 
    p.post_id,
    p.user_id,
    p.title,
    p.content,
    p.photo,
    p.video,
    p.created_at AS post_created_at,
    p.updated_at,
    f.created_at AS favorite_created_at,
    u.username AS user_name,
    u.avatar AS user_avatar,
    COALESCE(l.likes_count, 0) AS like_count,
    COALESCE(v.views_count, 0) AS views_count,
    COALESCE(c.comments_count, 0) AS comment_count,
    EXISTS (
        SELECT 1 FROM likes 
        WHERE likes.post_id = p.post_id 
        AND likes.user_id = f.user_id
    ) AS liked_by_user
FROM favorites f
JOIN posts p ON f.post_id = p.post_id
JOIN users u ON p.user_id = u.user_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS likes_count 
    FROM likes 
    GROUP BY post_id
) l ON p.post_id = l.post_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS views_count 
    FROM views 
    GROUP BY post_id
) v ON p.post_id = v.post_id
LEFT JOIN (
    SELECT post_id, COUNT(*) AS comments_count 
    FROM comments 
    GROUP BY post_id
) c ON p.post_id = c.post_id
WHERE f.user_id = $1
ORDER BY f.created_at DESC
LIMIT $2
OFFSET $3;

-- name: CheckFavoriteExists :one
SELECT EXISTS (
    SELECT 1 
    FROM favorites 
    WHERE user_id = $1 AND post_id = $2
) AS exists;