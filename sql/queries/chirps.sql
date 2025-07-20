-- name: CreateChirps :one
INSERT INTO chirps (body, user_id)
VALUES (
	$1,
	$2
	)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps 
ORDER BY (created_at) ASC;

-- name: GetAllChirpsByUser :many
SELECT * FROM chirps 
WHERE user_id = $1
ORDER BY (created_at) ASC;

-- name: FindChirp :one
SELECT * FROM chirps WHERE id = $1;



-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;
