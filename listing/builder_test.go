package listing_test

import (
	"errors"
	"testing"

	"marketplace/domain"
	"marketplace/listing"
)

func TestBuilder_BuildsValidListing(t *testing.T) {
	car := domain.NewCar("CAR-1", "Tesla", "Model 3", 2023, 15000, 5_000_000)
	seller := &domain.Seller{ID: "D-1", Type: domain.SellerDealer}

	l, err := listing.NewListingBuilder(car, seller).
		Title("Tesla Model 3").
		Description("Идеальное состояние").
		Photo("https://example.com/1.jpg").
		Photo("https://example.com/2.jpg").
		WithWarranty().
		WithTestDrive().
		Region("СПб").
		Commission(4).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Title != "Tesla Model 3" {
		t.Fatalf("title mismatch: %s", l.Title)
	}
	if len(l.Photos) != 2 {
		t.Fatalf("want 2 photos, got %d", len(l.Photos))
	}
	if l.Region != "СПб" {
		t.Fatalf("region mismatch: %s", l.Region)
	}
	if !l.HasWarranty || !l.HasTestDrive {
		t.Fatal("warranty and test drive must be true")
	}
	if l.Status != domain.StatusPending {
		t.Fatalf("want pending, got %s", l.Status)
	}
}

func TestBuilder_RequiresTitle(t *testing.T) {
	car := domain.NewCar("CAR-1", "Tesla", "Model 3", 2023, 15000, 5_000_000)
	seller := &domain.Seller{ID: "D-1", Type: domain.SellerDealer}

	_, err := listing.NewListingBuilder(car, seller).
		Description("Без заголовка").
		Build()

	if !errors.Is(err, listing.ErrMissingTitle) {
		t.Fatalf("want ErrMissingTitle, got %v", err)
	}
}

func TestBuilder_RejectsInvalidPrice(t *testing.T) {
	car := domain.NewCar("CAR-1", "Tesla", "Model 3", 2023, 15000, 0)
	seller := &domain.Seller{ID: "D-1", Type: domain.SellerDealer}

	_, err := listing.NewListingBuilder(car, seller).
		Title("Tesla").
		Build()

	if !errors.Is(err, listing.ErrInvalidPrice) {
		t.Fatalf("want ErrInvalidPrice, got %v", err)
	}
}
