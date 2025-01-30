package internal

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(connString string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	log.Println("INFO: successfully connected to the database")
	return conn, nil
}