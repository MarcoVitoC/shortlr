package internal

import (
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	port string
}

func NewServer(port string) *Config {
	return &Config{
		port: port,
	}
}

func (c *Config) Run() error {
	cache := NewRedisClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
	})

	service := Service{conn: cache}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shortlr", service.Generate)

	server := http.Server{
		Addr: c.port,
		Handler: mux,
	}

	log.Printf("INFO: server is running at %s", c.port)
	return server.ListenAndServe()
}