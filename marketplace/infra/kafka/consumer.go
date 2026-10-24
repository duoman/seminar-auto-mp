package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type DeliveryScheduledEvent struct {
	EventID    string `json:"event_id"`
	ListingID  string `json:"listing_id"`
	DeliveryID string `json:"delivery_id"`
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupID,
			StartOffset: kafka.FirstOffset, // ← добавь
		}),
	}
}
func (c *Consumer) Run(ctx context.Context, handler func(DeliveryScheduledEvent) error) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("[consumer] fetch: %v", err)
			continue
		}
		var e DeliveryScheduledEvent
		if err := json.Unmarshal(msg.Value, &e); err != nil {
			log.Printf("[consumer] unmarshal: %v", err)
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}
		if err := handler(e); err != nil {
			log.Printf("[consumer] handler: %v (will retry)", err)
			continue
		}
		_ = c.reader.CommitMessages(ctx, msg)
	}
}
