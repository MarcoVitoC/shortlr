package internal

import (
	"context"
	"crypto/rand"
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

const (
	length 	= 9
	charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
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

	WriteOK(w, data)
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
	if err := ReadJson(w, r, &payload); err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	shortlr, _ := s.repo.GetByLongUrl(ctx, payload.LongUrl)
	if shortlr != "" {
		WriteConflictResponse(w, "The long URL has been generated before.", nil)
		return
	}

	newShortlr, err := generateShortlr(ctx, s, payload)
	if err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	WriteOK(w, newShortlr)
}

func generateShortlr(ctx context.Context, s *Service, payload repository.Shortlr) (string, error) {
	if payload.LongUrl == "" {
		return "", errors.New("URL cannot be empty")
	}

	var newShortlr strings.Builder
	newShortlr.Grow(length)

	charsetLength := len(charset)
	for i:=0; i<length; i++ {
		k, _ := rand.Int(rand.Reader, big.NewInt(int64(charsetLength)))
		newShortlr.WriteByte(charset[k.Int64()])
	}

	return saveShortlr(ctx, s, payload, newShortlr.String())
}

func saveShortlr(ctx context.Context, s *Service, payload repository.Shortlr, newShortlr string) (string, error) {
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
	if err := ReadJson(w, r, &payload); err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	shortlr, err := s.repo.GetByLongUrl(ctx, payload.LongUrl)
	if shortlr != "" {
		WriteConflictResponse(w, "The long URL has been generated before.", err)
		return
	}

	updatedShortlr, err := s.repo.UpdateShortlr(ctx, repository.UpdateShortlrParams{
		LongUrl: payload.LongUrl,
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		ID: id,
	})

	if err != nil {
		WriteNotFoundResponse(w, err)
		return
	}

	expiration := time.Hour * 24 * 365
	s.cacheRepo.Set(ctx, updatedShortlr, payload.LongUrl, expiration)
	WriteOK(w, true)
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
	WriteOK(w, true)
}
