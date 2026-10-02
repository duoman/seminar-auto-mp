package listing

import "marketplace/domain"

type CommissionStrategy interface {
	Percent() int
}

type ModerationPolicy interface {
	RequiresManualReview() bool
}

type MarketplaceFactory interface {
	CreateListing(car *domain.Car, seller *domain.Seller) *domain.Listing
	CreateCommission() CommissionStrategy
	CreateModeration() ModerationPolicy
}

// --- Standard ---

type StandardFactory struct{}

func (StandardFactory) CreateListing(car *domain.Car, seller *domain.Seller) *domain.Listing {
	return FactoryFor(seller.Type).Create(car, seller)
}
func (StandardFactory) CreateCommission() CommissionStrategy { return standardCommission{} }
func (StandardFactory) CreateModeration() ModerationPolicy   { return standardModeration{} }

type standardCommission struct{}
func (standardCommission) Percent() int { return 3 }

type standardModeration struct{}
func (standardModeration) RequiresManualReview() bool { return false }

// --- Premium ---

type PremiumFactory struct{}

func (PremiumFactory) CreateListing(car *domain.Car, seller *domain.Seller) *domain.Listing {
	l := DealerListingFactory{}.Create(car, seller)
	l.Commission = 5
	return l
}
func (PremiumFactory) CreateCommission() CommissionStrategy { return premiumCommission{} }
func (PremiumFactory) CreateModeration() ModerationPolicy   { return premiumModeration{} }

type premiumCommission struct{}
func (premiumCommission) Percent() int { return 5 }

type premiumModeration struct{}
func (premiumModeration) RequiresManualReview() bool { return true }

// --- Auction ---

type AuctionFactory struct{}

func (AuctionFactory) CreateListing(car *domain.Car, seller *domain.Seller) *domain.Listing {
	return AuctionListingFactory{}.Create(car, seller)
}
func (AuctionFactory) CreateCommission() CommissionStrategy { return auctionCommission{} }
func (AuctionFactory) CreateModeration() ModerationPolicy   { return auctionModeration{} }

type auctionCommission struct{}
func (auctionCommission) Percent() int { return 7 }

type auctionModeration struct{}
func (auctionModeration) RequiresManualReview() bool { return true }
