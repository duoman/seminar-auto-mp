package ports

import "delivery/domain"

type DeliveryRepository interface {
	Save(d *domain.Delivery) error
	FindByID(id string) (*domain.Delivery, error)
	All() ([]*domain.Delivery, error)
}
