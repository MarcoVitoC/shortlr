package internal

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/MarcoVitoC/shortlr/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo *repository.Queries
	cacheRepo *redis.Client
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := s.repo.GetAllShortlr(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data == nil {
		data = []repository.Shortlr{}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) Generate(w http.ResponseWriter, r *http.Request) {
	var payload repository.Shortlr

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortlr, _ := s.repo.GetByLongUrl(context.Background(), payload.LongUrl)
	if shortlr != "" {
		if err := json.NewEncoder(w).Encode(shortlr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	newShortlr, err := generateShortlr(s, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(newShortlr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func generateShortlr(s *Service, payload repository.Shortlr) (string, error) {
	if payload.LongUrl == "" {
		return "", errors.New("ERROR: URL cannot be empty")
	}

	random := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	newShortlr := "short.lr/"
	for i:=0; i<9; i++ {
		k, _ := rand.Int(rand.Reader, big.NewInt(12))
		newShortlr += string(random[k.Int64()])
	}

	shortlr, err := s.repo.SaveShortlr(context.Background(), repository.SaveShortlrParams{
		ID: uuid.New(),
		LongUrl: payload.LongUrl,
		ShortUrl: newShortlr,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return "", errors.New(err.Error())
	}
	
	expiration := time.Hour * 24 * 365
	if err := s.cacheRepo.Set(context.Background(), newShortlr, payload.LongUrl, expiration); err != nil {
		return "", errors.New(err.Err().Error())
	}

	log.Printf("INFO: successfully generate new shortlr: %s", shortlr)
	return shortlr, nil
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	//
}
