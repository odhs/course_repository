package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"odhs/semana-tech-01-go-react-server-main/internal/api"
	"odhs/semana-tech-01-go-react-server-main/internal/store/pgstore"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	// Pool to PostgreSQL connections
	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("WS_DATABASE_USER"),
		os.Getenv("WS_DATABASE_PASSWORD"),
		os.Getenv("WS_DATABASE_HOST"),
		os.Getenv("WS_DATABASE_PORT"),
		os.Getenv("WS_DATABASE_NAME"),
	))

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	// Check if the database is online
	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	handler := api.NewHandler(pgstore.New(pool))

	go func() {
		if err := http.ListenAndServe(":8080", handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	slog.Info("Server started")

	// Create a channel to receive the signal Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

}
