package listing_test

import (
	"testing"

	"marketplace/domain"
	"marketplace/listing"
)

func TestClone_DeepCopiesPhotos(t *testing.T) {
	car := domain.NewCar("CAR-1", "Tesla", "Model 3", 2023, 15000, 5_000_000)
	seller := &domain.Seller{ID: "D-1", Type: domain.SellerDealer}

	original, _ := listing.NewListingBuilder(car, seller).
		Title("Оригинал").
		Photo("https://example.com/1.jpg").
		Build()

	clone := listing.Clone(original)
	clone.Photos = append(clone.Photos, "https://example.com/2.jpg")

	if len(original.Photos) != 1 {
		t.Fatalf("original photos changed: want 1, got %d", len(original.Photos))
	}
	if len(clone.Photos) != 2 {
		t.Fatalf("clone photos: want 2, got %d", len(clone.Photos))
	}
}

func TestClone_DeepCopiesCar(t *testing.T) {
	car := domain.NewCar("CAR-1", "Tesla", "Model 3", 2023, 15000, 5_000_000)
	seller := &domain.Seller{ID: "D-1", Type: domain.SellerDealer}

	original, _ := listing.NewListingBuilder(car, seller).Title("Оригинал").Build()

	clone := listing.Clone(original)
	clone.Car.Price = 1

	if original.Car.Price == 1 {
		t.Fatal("changing clone.Car.Price must not affect original")
	}
}

func TestClone_MarksAsDraftCopy(t *testing.T) {
	car := domain.NewCar("CAR-1", "Tesla", "Model 3", 2023, 15000, 5_000_000)
	seller := &domain.Seller{ID: "D-1", Type: domain.SellerDealer}

	original, _ := listing.NewListingBuilder(car, seller).Title("Оригинал").Build()
	clone := listing.Clone(original)

	if clone.Status != domain.StatusDraft {
		t.Fatalf("clone must be draft, got %s", clone.Status)
	}
	if clone.ID == original.ID {
		t.Fatal("clone ID must differ from original")
	}
}
