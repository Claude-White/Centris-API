package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"

	"centris-api/internal/repository"
)

type Server struct {
	port    int
	db      *pgx.Conn
	queries *repository.Queries
}

func CreateServer(db *pgx.Conn) Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	repo := repository.New(db)
	NewServer := &Server{
		port:    port,
		db:      db,
		queries: repo,
	}

	return *NewServer
}

func NewServer(db *pgx.Conn) *http.Server {
	NewServer := CreateServer(db)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
