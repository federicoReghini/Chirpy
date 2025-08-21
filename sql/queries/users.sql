-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING id, created_at, updated_at, email;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE users.email = $1;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users 
WHERE users.id = (
  SELECT user_id FROM refresh_tokens
  WHERE refresh_tokens.token = $1
);

-- name: UpdateUser :one
UPDATE users
  SET email = $1, hashed_password = $2
  WHERE users.id = $3
  RETURNING *;
