package listing

import (
	"marketplace/domain/catalog"
	"marketplace/domain/sales"
)

type CommissionStrategy interface{ Percent() int }
type ModerationPolicy interface{ Manual() bool }

type MarketplaceFactory interface {
	CreateListing(id sales.ListingID, car *catalog.Car, seller *sales.Seller) (*sales.Listing, error)
	CreateCommission() CommissionStrategy
	CreateModeration() ModerationPolicy
}

type StandardFactory struct{}
func (StandardFactory) CreateListing(id sales.ListingID, car *catalog.Car, seller *sales.Seller) (*sales.Listing, error) {
	return FactoryFor(seller.Type()).Create(id, car, seller)
}
func (StandardFactory) CreateCommission() CommissionStrategy { return stdCommission{} }
func (StandardFactory) CreateModeration() ModerationPolicy   { return stdModeration{} }
type stdCommission struct{}
func (stdCommission) Percent() int { return 3 }
type stdModeration struct{}
func (stdModeration) Manual() bool { return false }

type PremiumFactory struct{}
func (PremiumFactory) CreateListing(id sales.ListingID, car *catalog.Car, seller *sales.Seller) (*sales.Listing, error) {
	l, err := DealerListingFactory{}.Create(id, car, seller)
	if err != nil { return nil, err }
	c, _ := sales.NewDiscountPercent(5)
	_ = l.SetCommission(c)
	return l, nil
}
func (PremiumFactory) CreateCommission() CommissionStrategy { return premCommission{} }
func (PremiumFactory) CreateModeration() ModerationPolicy   { return premModeration{} }
type premCommission struct{}
func (premCommission) Percent() int { return 5 }
type premModeration struct{}
func (premModeration) Manual() bool { return true }

type AuctionFactory struct{}
func (AuctionFactory) CreateListing(id sales.ListingID, car *catalog.Car, seller *sales.Seller) (*sales.Listing, error) {
	return AuctionListingFactory{}.Create(id, car, seller)
}
func (AuctionFactory) CreateCommission() CommissionStrategy { return aucCommission{} }
func (AuctionFactory) CreateModeration() ModerationPolicy   { return aucModeration{} }
type aucCommission struct{}
func (aucCommission) Percent() int { return 7 }
type aucModeration struct{}
func (aucModeration) Manual() bool { return true }
