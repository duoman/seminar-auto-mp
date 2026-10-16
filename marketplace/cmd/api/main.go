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
	delivery := pickDeliveryClient()

	server := httpapi.NewServer(repo, delivery)
	addr := ":" + envOr("PORT", "8080")
	log.Printf("marketplace on %s (storage=%s, delivery=%s)",
		addr, envOr("STORAGE", "memory"), envOr("DELIVERY_URL", "http://localhost:8081"))

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
		return infra.NewPostgresListingRepository(db)
	default:
		return infra.NewInMemoryListingRepository()
	}
}

func pickDeliveryClient() ports.DeliveryClient {
	url := envOr("DELIVERY_URL", "http://localhost:8081")
	return infra.NewHTTPDeliveryClient(url)
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
