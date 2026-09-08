package main

import (
	"fmt"
	"log"
	"os"

	"github.com/durgasatkar/olx-api/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate <up|down>")
	}
	cfg := config.MustLoad()
	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseUrl,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatalf("failed to migrate up: %v", err)
		}
		fmt.Println("Migrating up")
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("failed to migrate down: %v", err)
		}
		fmt.Println("Migrating down")
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}

	fmt.Println("running migration")
}
