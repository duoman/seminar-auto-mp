package listing_test

import (
	"testing"

	"marketplace/listing"
)

func TestMarketplaceFamilies(t *testing.T) {
	tests := []struct {
		name           string
		factory        listing.MarketplaceFactory
		wantCommission int
		wantReview     bool
	}{
		{"standard", listing.StandardFactory{}, 3, false},
		{"premium", listing.PremiumFactory{}, 5, true},
		{"auction", listing.AuctionFactory{}, 7, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comm := tt.factory.CreateCommission().Percent()
			mod := tt.factory.CreateModeration().RequiresManualReview()

			if comm != tt.wantCommission {
				t.Fatalf("commission: want %d, got %d", tt.wantCommission, comm)
			}
			if mod != tt.wantReview {
				t.Fatalf("review: want %v, got %v", tt.wantReview, mod)
			}
		})
	}
}
