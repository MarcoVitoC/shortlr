package internal

import (
	"log"
	"net/http"
	"os"

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

	db, err := InitDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("ERROR: failed to connect to database")
	}
	defer db.Close()

	redis := NewRedisClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
	})

	service := Service{conn: db, cache: redis}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shortlr", service.Generate)

	server := http.Server{
		Addr: c.port,
		Handler: mux,
	}

	log.Printf("INFO: server is running at %s", c.port)
	return server.ListenAndServe()
}