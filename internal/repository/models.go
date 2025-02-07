// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Shortlr struct {
	ID          uuid.UUID        `json:"id"`
	LongUrl     string           `json:"long_url"`
	ShortUrl    string           `json:"short_url"`
	AccessCount pgtype.Int8      `json:"access_count"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}
