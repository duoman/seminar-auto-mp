package main

import (
	"fmt"

	"marketplace/domain"
	"marketplace/infra"
	"marketplace/listing"
)

func main() {
	fmt.Println("=== CLI-демонстрация паттернов ===")

	car := domain.NewCar("CAR-1", "Tesla", "Model 3", 2023, 15000, 5_000_000)

	// Factory Method
	fmt.Println("\n-- Factory Method --")
	for _, st := range []domain.SellerType{domain.SellerDealer, domain.SellerPrivate, domain.SellerAuction} {
		seller := &domain.Seller{ID: "S-1", Name: "Seller", Type: st}
		l := listing.FactoryFor(st).Create(car, seller)
		fmt.Printf("[%-8s] %s (комиссия %d%%)\n", st, l.Title, l.Commission)
	}

	// Builder
	fmt.Println("\n-- Builder --")
	l, err := listing.NewListingBuilder(car, &domain.Seller{ID: "D-1", Type: domain.SellerDealer}).
		Title("Tesla Model 3, идеальное состояние").
		Description("Один владелец").
		Photo("https://example.com/1.jpg").
		Photo("https://example.com/2.jpg").
		WithWarranty().WithTestDrive().
		Region("СПб").Commission(4).
		Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s | фото: %d | регион: %s\n", l.Title, len(l.Photos), l.Region)

	// Prototype
	fmt.Println("\n-- Prototype --")
	clone := listing.Clone(l)
	clone.Photos = append(clone.Photos, "https://example.com/3.jpg")
	fmt.Printf("original фото: %d\n", len(l.Photos))
	fmt.Printf("clone фото:    %d\n", len(clone.Photos))
	fmt.Println("Оригинал не изменился:", len(l.Photos) == 2)

	// Abstract Factory
	fmt.Println("\n-- Abstract Factory --")
	markets := map[string]listing.MarketplaceFactory{
		"standard": listing.StandardFactory{},
		"premium":  listing.PremiumFactory{},
		"auction":  listing.AuctionFactory{},
	}
	for name, mf := range markets {
		fmt.Printf("[%-8s] комиссия %d%%, модерация: %v\n",
			name, mf.CreateCommission().Percent(), mf.CreateModeration().RequiresManualReview())
	}

	// Repository
	repo := infra.NewInMemoryListingRepository()
	_ = repo.Save(l)
	_ = repo.Save(clone)
	all, _ := repo.All()
	fmt.Println("\nСохранено:", len(all))
}
