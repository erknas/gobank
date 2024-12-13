package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env file", err)
	}

	var (
		listenAddr = os.Getenv("LISTEN_ADDR")
		connStr    = os.Getenv("DATABASE_URL")
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	store, err := NewStorage(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close(ctx)

	logger := NewLogger(store)

	srv := NewServer(listenAddr, logger)
	srv.Run(ctx)
}
