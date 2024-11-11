package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env file", err)
	}

	var (
		listenAddr = os.Getenv("LISTEN_ADDR")
		connStr    = os.Getenv("DATABASE_URL")
		ctx        = context.Background()
	)

	store, err := NewStorage(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close(ctx)

	srv := NewServer(listenAddr, store)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
