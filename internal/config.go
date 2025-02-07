package internal

import (
	"log"
	"net/http"
	"os"

	"github.com/MarcoVitoC/shortlr/internal/repository"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR: failed to load .env file")
	}

	conn, err := InitDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("ERROR: failed to connect to database")
	}
	defer conn.Close()

	redisClient := NewRedisClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
	})

	repo := repository.New(conn)
	service := Service{
		repo: repo, 
		cacheRepo: redisClient,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /shortlr", service.GetAll)
	mux.HandleFunc("POST /shortlr", service.Generate)
	mux.HandleFunc("DELETE /shortlr/{id}", service.Delete)

	server := http.Server{
		Addr: c.port,
		Handler: mux,
	}

	log.Printf("INFO: server is running at %s", c.port)
	return server.ListenAndServe()
}