package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		connStr        = os.Getenv("DATABASE_URL")
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
	)

	m, err := migrate.New(migrationsPath, connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no new migrations")
			return
		}
		log.Fatal(err)
	}

	fmt.Println("successful migration")
}
