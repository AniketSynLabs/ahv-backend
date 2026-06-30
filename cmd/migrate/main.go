package main

import (
	"log"

	"ahv-worldwide/config"
	"ahv-worldwide/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	if err := db.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("✅ Migration complete")
}
