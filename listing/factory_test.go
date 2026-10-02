package listing_test

import (
	"strings"
	"testing"

	"marketplace/domain"
	"marketplace/listing"
)

func TestDealerListingFactory(t *testing.T) {
	car := domain.NewCar("CAR-1", "Tesla", "Model 3", 2023, 15000, 5_000_000)
	seller := &domain.Seller{ID: "D-1", Name: "AutoDealer", Type: domain.SellerDealer}

	l := listing.DealerListingFactory{}.Create(car, seller)

	if !l.HasWarranty {
		t.Fatal("dealer listing must have warranty")
	}
	if l.Commission != 3 {
		t.Fatalf("want commission 3, got %d", l.Commission)
	}
	if l.Status != domain.StatusPending {
		t.Fatalf("want pending, got %s", l.Status)
	}
}

func TestPrivateListingFactory(t *testing.T) {
	car := domain.NewCar("CAR-1", "Lada", "Vesta", 2020, 80000, 800_000)
	seller := &domain.Seller{ID: "P-1", Type: domain.SellerPrivate}

	l := listing.PrivateListingFactory{}.Create(car, seller)

	if l.HasWarranty {
		t.Fatal("private listing must not have warranty")
	}
	if l.Commission != 1 {
		t.Fatalf("want commission 1, got %d", l.Commission)
	}
	if !l.HasTestDrive {
		t.Fatal("private listing must have test drive")
	}
}

func TestAuctionListingFactory(t *testing.T) {
	car := domain.NewCar("CAR-1", "BMW", "X5", 2021, 40000, 4_000_000)
	seller := &domain.Seller{ID: "A-1", Type: domain.SellerAuction}

	l := listing.AuctionListingFactory{}.Create(car, seller)

	if l.Status != domain.StatusDraft {
		t.Fatalf("auction listing must be draft, got %s", l.Status)
	}
	if l.HasWarranty || l.HasTestDrive {
		t.Fatal("auction listing must not have warranty or test drive")
	}
}

func TestFactoryFor_SelectsByType(t *testing.T) {
	car := domain.NewCar("CAR-1", "Tesla", "Model 3", 2023, 15000, 5_000_000)

	tests := []struct {
		name       string
		sellerType domain.SellerType
		wantSubstr string
	}{
		{"dealer", domain.SellerDealer, "официального дилера"},
		{"private", domain.SellerPrivate, "пробег"},
		{"auction", domain.SellerAuction, "Аукцион"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seller := &domain.Seller{ID: "S-1", Type: tt.sellerType}
			factory := listing.FactoryFor(tt.sellerType)
			l := factory.Create(car, seller)

			if !strings.Contains(l.Title, tt.wantSubstr) {
				t.Fatalf("title %q must contain %q", l.Title, tt.wantSubstr)
			}
		})
	}
}
