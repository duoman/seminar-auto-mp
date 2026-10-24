package sales

// OutboxEvent — событие для публикации в Kafka, хранится в той же транзакции, что и Listing.
type OutboxEvent struct {
	ID      string
	Type    string // "ListingSold"
	Payload []byte
}

const EventListingSold = "ListingSold"
