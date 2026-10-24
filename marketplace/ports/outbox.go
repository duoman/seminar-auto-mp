package ports

import "marketplace/domain/sales"

type OutboxRelay interface {
	Run()
}

type UnpublishedEvents interface {
	FetchUnpublished(limit int) ([]sales.OutboxEvent, error)
	MarkPublished(id string) error
}
