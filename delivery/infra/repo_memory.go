package infra

import (
	"errors"
	"sync"

	"delivery/domain"
)

type InMemoryDeliveryRepository struct {
	mu         sync.RWMutex
	deliveries map[string]*domain.Delivery
}

func NewInMemoryDeliveryRepository() *InMemoryDeliveryRepository {
	return &InMemoryDeliveryRepository{deliveries: map[string]*domain.Delivery{}}
}

func (r *InMemoryDeliveryRepository) Save(d *domain.Delivery) error {
	if d == nil {
		return errors.New("nil delivery")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deliveries[d.ID] = d
	return nil
}

func (r *InMemoryDeliveryRepository) FindByID(id string) (*domain.Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.deliveries[id]
	if !ok {
		return nil, errors.New("delivery not found")
	}
	return d, nil
}

func (r *InMemoryDeliveryRepository) All() ([]*domain.Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Delivery, 0, len(r.deliveries))
	for _, d := range r.deliveries {
		out = append(out, d)
	}
	return out, nil
}
