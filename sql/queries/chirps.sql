-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY 
CASE  WHEN $1 = 'asc' THEN  created_at END ASC,
CASE  WHEN $1 = 'desc' THEN  created_at END DESC;

-- name: GetChirpsByUser :many
SELECT *
  FROM chirps 
  WHERE chirps.user_id = $1;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE chirps.id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps 
  WHERE chirps.user_id = $1 AND chirps.id = $2;

