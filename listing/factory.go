package listing

import (
	"fmt"

	"marketplace/domain"
)

type ListingFactory interface {
	Create(car *domain.Car, seller *domain.Seller) *domain.Listing
}

type DealerListingFactory struct{}

func (DealerListingFactory) Create(car *domain.Car, seller *domain.Seller) *domain.Listing {
	return &domain.Listing{
		ID:           fmt.Sprintf("L-%s-%s", seller.ID, car.ID),
		Car:          car, Seller: seller,
		Title:        fmt.Sprintf("%s %s %d, от официального дилера", car.Brand, car.Model, car.Year),
		Description:  "Гарантия, сервисная история, проверка по 120 пунктам.",
		Price:        car.Price, Status: domain.StatusPending,
		HasWarranty:  true, HasTestDrive: true, Region: "Moscow", Commission: 3,
	}
}

type PrivateListingFactory struct{}

func (PrivateListingFactory) Create(car *domain.Car, seller *domain.Seller) *domain.Listing {
	return &domain.Listing{
		ID:           fmt.Sprintf("L-%s-%s", seller.ID, car.ID),
		Car:          car, Seller: seller,
		Title:        fmt.Sprintf("%s %s %d, пробег %d км", car.Brand, car.Model, car.Year, car.Mileage),
		Description:  "Продаю личную машину. Обслуживание по регламенту.",
		Price:        car.Price, Status: domain.StatusPending,
		HasWarranty:  false, HasTestDrive: true, Region: "Moscow", Commission: 1,
	}
}

type AuctionListingFactory struct{}

func (AuctionListingFactory) Create(car *domain.Car, seller *domain.Seller) *domain.Listing {
	return &domain.Listing{
		ID:           fmt.Sprintf("L-%s-%s", seller.ID, car.ID),
		Car:          car, Seller: seller,
		Title:        fmt.Sprintf("Аукцион: %s %s %d", car.Brand, car.Model, car.Year),
		Description:  "Стартовая цена ниже рыночной. Торги 24 часа.",
		Price:        car.Price, Status: domain.StatusDraft,
		HasWarranty:  false, HasTestDrive: false, Region: "Moscow", Commission: 5,
	}
}

func FactoryFor(sellerType domain.SellerType) ListingFactory {
	switch sellerType {
	case domain.SellerDealer:
		return DealerListingFactory{}
	case domain.SellerPrivate:
		return PrivateListingFactory{}
	case domain.SellerAuction:
		return AuctionListingFactory{}
	default:
		panic("unknown seller type: " + string(sellerType))
	}
}
