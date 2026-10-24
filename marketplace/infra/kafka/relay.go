package kafka

import (
	"context"
	"log"
	"time"

	"marketplace/ports"

	"github.com/segmentio/kafka-go"
)

type Relay struct {
	repo   ports.UnpublishedEvents
	writer *kafka.Writer
	topic  string
}

func NewRelay(repo ports.UnpublishedEvents, brokers []string, topic string) *Relay {
	return &Relay{
		repo: repo,
		topic: topic,
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (r *Relay) Run(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			_ = r.writer.Close()
			return
		case <-ticker.C:
			events, err := r.repo.FetchUnpublished(50)
			if err != nil { log.Printf("[relay] fetch: %v", err); continue }
			for _, e := range events {
				if err := r.writer.WriteMessages(ctx, kafka.Message{
					Key:   []byte(e.ID),
					Value: e.Payload,
				}); err != nil {
					log.Printf("[relay] write: %v", err)
					continue
				}
				_ = r.repo.MarkPublished(e.ID)
				log.Printf("[relay] published %s (%s)", e.ID, e.Type)
			}
		}
	}
}
