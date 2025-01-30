package internal

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	conn *pgxpool.Pool
	cache *redis.Client
}

func (s *Service) Generate(w http.ResponseWriter, r *http.Request) {
	var payload Shortlr

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//
}

