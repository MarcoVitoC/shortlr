package internal

import (
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	conn *redis.Client
}

func (s *Service) Generate(w http.ResponseWriter, r *http.Request) {
	//
}