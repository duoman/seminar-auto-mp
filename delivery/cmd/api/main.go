package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"delivery/behavior"
	"delivery/httpapi"
	"delivery/infra"
	kafkainfra "delivery/infra/kafka"
	"delivery/structural"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", envOr("DATABASE_URL", ""))
	if err != nil { log.Fatal(err) }
	if err := db.Ping(); err != nil { log.Fatal(err) }
	defer db.Close()

	repo := infra.NewPostgresRepo(db)
	outbox := infra.NewOutboxRepo(db)
	brokers := strings.Split(envOr("KAFKA_BROKERS", "kafka:9092"), ",")
	ctx := context.Background()

	pub := &behavior.Publisher{}
	pub.Subscribe(behavior.LogObserver{Name: "customer"})
	pub.Subscribe(behavior.LogObserver{Name: "warehouse"})
	pub.Subscribe(behavior.LogObserver{Name: "analytics"})

	tpl := structural.NewTemplates()
	facade := structural.NewScheduleFacade(repo, tpl, pub)

	go kafkainfra.NewRelay(outbox, brokers, "delivery-events").Run(ctx)

	consumer := kafkainfra.NewConsumer(brokers, "listing-events", "delivery-service")
	go consumer.Run(ctx, func(e kafkainfra.ListingSoldEvent) error {
		id := "DEL-" + e.ListingID
		if _, err := repo.FindByID(id); err == nil {
			log.Printf("[consumer] delivery %s already exists", id)
			return nil
		}
		d, err := facade.Schedule(id, e.ListingID, e.CustomerID, e.Type, e.FromCity, e.ToCity)
		if err != nil { return err }
		payload, _ := json.Marshal(map[string]any{
			"event_id":    fmt.Sprintf("evt-del-%s-%d", d.ID(), time.Now().UnixNano()),
			"listing_id":  d.ListingID(),
			"delivery_id": d.ID(),
		})
		outID := fmt.Sprintf("outbox-del-%s-%d", d.ID(), time.Now().UnixNano())
		if err := repo.SaveWithOutbox(d, outID, "DeliveryScheduled", payload); err != nil { return err }
		log.Printf("[consumer] created delivery %s for listing %s", d.ID(), d.ListingID())
		return nil
	})

	server := httpapi.NewServer(repo, facade, pub)
	addr := ":" + envOr("PORT", "8081")
	log.Printf("delivery on %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Routes()))
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}
