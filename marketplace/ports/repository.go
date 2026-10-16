package ports

import "marketplace/domain"

type ListingRepository interface {
	Save(l *domain.Listing) error
	FindByID(id string) (*domain.Listing, error)
	All() ([]*domain.Listing, error)
}
