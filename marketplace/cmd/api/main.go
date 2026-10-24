package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"marketplace/domain/sales"
	"marketplace/httpapi"
	"marketplace/infra"
	kafkainfra "marketplace/infra/kafka"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", envOr("DATABASE_URL", ""))
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := infra.NewPostgresRepo(db)
	outbox := infra.NewOutboxRepo(db)

	brokers := strings.Split(envOr("KAFKA_BROKERS", "kafka:9092"), ",")
	ctx := context.Background()

	// Relay: публикует события из outbox в Kafka
	relay := kafkainfra.NewRelay(outbox, brokers, "listing-events")
	go relay.Run(ctx)

	// Consumer: слушает delivery-events, обновляет listings.delivery_id
	consumer := kafkainfra.NewConsumer(brokers, "delivery-events", "marketplace-service")
	go consumer.Run(ctx, func(e kafkainfra.DeliveryScheduledEvent) error {
		id, _ := sales.NewListingID(e.ListingID)
		l, err := repo.FindByID(id)
		if err != nil {
			log.Printf("[consumer] listing %s not found: %v", e.ListingID, err)
			return nil
		}
		if err := l.AttachDelivery(e.DeliveryID); err != nil {
			log.Printf("[consumer] attach: %v", err)
			return nil
		}
		if err := repo.Save(l); err != nil {
			return err
		}
		log.Printf("[consumer] attached delivery %s to listing %s", e.DeliveryID, e.ListingID)
		return nil
	})

	server := httpapi.NewServer(repo)
	addr := ":" + envOr("PORT", "8080")
	log.Printf("marketplace on %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Routes()))
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
