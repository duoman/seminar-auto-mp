package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"marketplace/httpapi"
	"marketplace/infra"
	"marketplace/ports"

	_ "github.com/lib/pq"
)

func main() {
	repo := pickRepository()
	defer closeRepo(repo)

	server := httpapi.NewServer(repo)
	addr := ":" + envOr("PORT", "8080")
	log.Printf("listening on %s (storage=%s)", addr, envOr("STORAGE", "memory"))

	if err := http.ListenAndServe(addr, server.Routes()); err != nil {
		log.Fatal(err)
	}
}

func pickRepository() ports.ListingRepository {
	switch envOr("STORAGE", "memory") {
	case "postgres":
		db, err := sql.Open("postgres", envOr("DATABASE_URL", ""))
		if err != nil {
			log.Fatalf("open db: %v", err)
		}
		if err := db.Ping(); err != nil {
			log.Fatalf("ping db: %v", err)
		}
		log.Println("using Postgres storage")
		return infra.NewPostgresListingRepository(db)
	default:
		log.Println("using in-memory storage")
		return infra.NewInMemoryListingRepository()
	}
}

func closeRepo(repo ports.ListingRepository) {
	if pg, ok := repo.(*infra.PostgresListingRepository); ok {
		_ = pg.Close()
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
