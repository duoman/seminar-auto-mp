package ports

import "marketplace/domain/sales"

type ListingRepository interface {
	Save(l *sales.Listing) error
	SaveWithOutbox(l *sales.Listing, evt sales.OutboxEvent) error
	FindByID(id sales.ListingID) (*sales.Listing, error)
	All() ([]*sales.Listing, error)
	AttachDelivery(listingID sales.ListingID, deliveryID string) error
}
