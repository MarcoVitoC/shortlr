package internal

import (
	"encoding/json"
	"net/http"

	"github.com/MarcoVitoC/shortlr/internal/repository"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo *repository.Queries
	cacheRepo *redis.Client
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

