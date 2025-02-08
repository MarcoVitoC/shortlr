-- name: GetAllShortlr :many
SELECT * FROM shortlrs;

-- name: GetByLongUrl :one
SELECT short_url FROM shortlrs
WHERE long_url = $1;

-- name: SaveShortlr :one
INSERT INTO shortlrs (
    id, long_url, short_url, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING short_url;

-- name: UpdateShortlr :one
UPDATE shortlrs
SET long_url = $1, updated_at = $2
WHERE id = $3
RETURNING short_url;

-- name: IncrementAccessCount :exec
UPDATE shortlrs
SET access_count = access_count + 1
WHERE short_url = $1;

-- name: DeleteShortlr :one
DELETE FROM shortlrs
WHERE id = $1
RETURNING short_url;
