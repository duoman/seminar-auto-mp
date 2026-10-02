package listing

import (
	"errors"

	"marketplace/domain"
)

var (
	ErrMissingTitle = errors.New("listing: title is required")
	ErrInvalidPrice = errors.New("listing: price must be positive")
)

type ListingBuilder struct {
	listing domain.Listing
}

func NewListingBuilder(car *domain.Car, seller *domain.Seller) *ListingBuilder {
	return &ListingBuilder{
		listing: domain.Listing{
			ID: "L-" + seller.ID + "-" + car.ID,
			Car: car, Seller: seller,
			Price: car.Price, Status: domain.StatusDraft, Region: "Moscow",
		},
	}
}

func (b *ListingBuilder) Title(v string) *ListingBuilder       { b.listing.Title = v; return b }
func (b *ListingBuilder) Description(v string) *ListingBuilder { b.listing.Description = v; return b }
func (b *ListingBuilder) Photo(url string) *ListingBuilder     { b.listing.Photos = append(b.listing.Photos, url); return b }
func (b *ListingBuilder) WithWarranty() *ListingBuilder        { b.listing.HasWarranty = true; return b }
func (b *ListingBuilder) WithTestDrive() *ListingBuilder       { b.listing.HasTestDrive = true; return b }
func (b *ListingBuilder) Region(v string) *ListingBuilder      { b.listing.Region = v; return b }
func (b *ListingBuilder) Commission(v int) *ListingBuilder     { b.listing.Commission = v; return b }

func (b *ListingBuilder) Build() (*domain.Listing, error) {
	if b.listing.Title == "" {
		return nil, ErrMissingTitle
	}
	if b.listing.Price <= 0 {
		return nil, ErrInvalidPrice
	}
	b.listing.Status = domain.StatusPending
	return &b.listing, nil
}
