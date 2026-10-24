package ports

import "delivery/domain"

type Repository interface {
	Save(d *domain.Delivery) error
	SaveWithOutbox(d *domain.Delivery, eventID, eventType string, payload []byte) error
	FindByID(id string) (*domain.Delivery, error)
	All() ([]*domain.Delivery, error)
}

type OutboxEvent struct {
	ID, Type string
	Payload  []byte
}

type OutboxRepo interface {
	FetchUnpublished(limit int) ([]OutboxEvent, error)
	MarkPublished(id string) error
}
