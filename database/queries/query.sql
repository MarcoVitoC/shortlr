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
