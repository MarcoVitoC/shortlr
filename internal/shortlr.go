package internal

import (
	"time"

	"github.com/google/uuid"
)

type Shortlr struct {
	ID 			uuid.UUID `json:"id"`
	LongUrl  	string    `json:"long_url"`
	ShortUrl 	string    `json:"short_url"`
	AccessCount int64	  `json:"access_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}