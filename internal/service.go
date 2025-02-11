package internal

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"net/http"
	"strings"
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
	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	data, err := s.repo.GetAllShortlr(ctx)
	if err != nil {
		WriteInternalServerErrorResponse(w, err)
		return
	}

	WriteOKResponse(w, data)
}

func (s *Service) Redirect(w http.ResponseWriter, r *http.Request) {
	urlPrefix := "https://"

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	key := r.PathValue("shortlr")
	url, err := s.cacheRepo.Get(ctx, key).Result()
	if err != nil {
		WriteInternalServerErrorResponse(w, err)
		return
	}

	if !strings.HasPrefix(url, urlPrefix) {
		url = urlPrefix + url
	}

	if err := s.repo.IncrementAccessCount(ctx, key); err != nil {
		WriteInternalServerErrorResponse(w, err)
		return
	}
	
	http.Redirect(w, r, url, http.StatusFound)
}

func (s *Service) Generate(w http.ResponseWriter, r *http.Request) {
	var payload repository.Shortlr
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	shortlr, _ := s.repo.GetByLongUrl(ctx, payload.LongUrl)
	if shortlr != "" {
		WriteOKResponse(w, shortlr)
		return
	}

	newShortlr, err := generateShortlr(ctx, s, payload)
	if err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	WriteOKResponse(w, newShortlr)
}

func generateShortlr(ctx context.Context, s *Service, payload repository.Shortlr) (string, error) {
	if payload.LongUrl == "" {
		return "", errors.New("ERROR: URL cannot be empty")
	}

	random, newShortlr := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", ""
	for i:=0; i<9; i++ {
		k, _ := rand.Int(rand.Reader, big.NewInt(12))
		newShortlr += string(random[k.Int64()])
	}

	shortlr, err := s.repo.SaveShortlr(ctx, repository.SaveShortlrParams{
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
	s.cacheRepo.Set(ctx, newShortlr, payload.LongUrl, expiration)

	log.Printf("INFO: successfully generate new shortlr: %s", shortlr)
	return shortlr, nil
}

func (s *Service) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	var payload repository.Shortlr
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	shortlr, err := s.repo.GetByLongUrl(ctx, payload.LongUrl)
	if shortlr != "" {
		WriteConflictResponse(w, err)
		return
	}

	updatedShortlr, err := s.repo.UpdateShortlr(ctx, repository.UpdateShortlrParams{
		LongUrl: payload.LongUrl,
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		ID: id,
	})

	if err != nil {
		WriteConflictResponse(w, err)
		return
	}

	expiration := time.Hour * 24 * 365
	s.cacheRepo.Set(ctx, updatedShortlr, payload.LongUrl, expiration)
	WriteOKResponse(w, updatedShortlr)
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	shortlr, err := s.repo.DeleteShortlr(ctx, id)
	if err != nil {
		WriteNotFoundResponse(w, err)
		return
	}

	s.cacheRepo.Del(ctx, shortlr)
	WriteOKResponse(w, true)
}
