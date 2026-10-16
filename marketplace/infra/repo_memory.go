package infra

import (
	"errors"
	"sync"

	"marketplace/domain"
)

type InMemoryListingRepository struct {
	mu       sync.RWMutex
	listings map[string]*domain.Listing
}

func NewInMemoryListingRepository() *InMemoryListingRepository {
	return &InMemoryListingRepository{listings: map[string]*domain.Listing{}}
}

func (r *InMemoryListingRepository) Save(l *domain.Listing) error {
	if l == nil {
		return errors.New("nil listing")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listings[l.ID] = l
	return nil
}

func (r *InMemoryListingRepository) FindByID(id string) (*domain.Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.listings[id]
	if !ok {
		return nil, errors.New("listing not found")
	}
	return l, nil
}

func (r *InMemoryListingRepository) All() ([]*domain.Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Listing, 0, len(r.listings))
	for _, l := range r.listings {
		out = append(out, l)
	}
	return out, nil
}
