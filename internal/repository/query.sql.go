// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const deleteShortlr = `-- name: DeleteShortlr :one
DELETE FROM shortlrs
WHERE id = $1
RETURNING short_url
`

func (q *Queries) DeleteShortlr(ctx context.Context, id uuid.UUID) (string, error) {
	row := q.db.QueryRow(ctx, deleteShortlr, id)
	var short_url string
	err := row.Scan(&short_url)
	return short_url, err
}

const getAllShortlr = `-- name: GetAllShortlr :many
SELECT id, long_url, short_url, access_count, created_at, updated_at FROM shortlrs
`

func (q *Queries) GetAllShortlr(ctx context.Context) ([]Shortlr, error) {
	rows, err := q.db.Query(ctx, getAllShortlr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Shortlr
	for rows.Next() {
		var i Shortlr
		if err := rows.Scan(
			&i.ID,
			&i.LongUrl,
			&i.ShortUrl,
			&i.AccessCount,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getByLongUrl = `-- name: GetByLongUrl :one
SELECT short_url FROM shortlrs
WHERE long_url = $1
`

func (q *Queries) GetByLongUrl(ctx context.Context, longUrl string) (string, error) {
	row := q.db.QueryRow(ctx, getByLongUrl, longUrl)
	var short_url string
	err := row.Scan(&short_url)
	return short_url, err
}

const incrementAccessCount = `-- name: IncrementAccessCount :exec
UPDATE shortlrs
SET access_count = access_count + 1
WHERE short_url = $1
`

func (q *Queries) IncrementAccessCount(ctx context.Context, shortUrl string) error {
	_, err := q.db.Exec(ctx, incrementAccessCount, shortUrl)
	return err
}

const saveShortlr = `-- name: SaveShortlr :one
INSERT INTO shortlrs (
    id, long_url, short_url, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING short_url
`

type SaveShortlrParams struct {
	ID        uuid.UUID        `json:"id"`
	LongUrl   string           `json:"long_url"`
	ShortUrl  string           `json:"short_url"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) SaveShortlr(ctx context.Context, arg SaveShortlrParams) (string, error) {
	row := q.db.QueryRow(ctx, saveShortlr,
		arg.ID,
		arg.LongUrl,
		arg.ShortUrl,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var short_url string
	err := row.Scan(&short_url)
	return short_url, err
}

const updateShortlr = `-- name: UpdateShortlr :one
UPDATE shortlrs
SET long_url = $1, updated_at = $2
WHERE id = $3
RETURNING short_url
`

type UpdateShortlrParams struct {
	LongUrl   string           `json:"long_url"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	ID        uuid.UUID        `json:"id"`
}

func (q *Queries) UpdateShortlr(ctx context.Context, arg UpdateShortlrParams) (string, error) {
	row := q.db.QueryRow(ctx, updateShortlr, arg.LongUrl, arg.UpdatedAt, arg.ID)
	var short_url string
	err := row.Scan(&short_url)
	return short_url, err
}
