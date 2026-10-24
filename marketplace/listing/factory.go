package listing

import (
	"fmt"

	"marketplace/domain/catalog"
	"marketplace/domain/sales"
)

type ListingFactory interface {
	Create(id sales.ListingID, car *catalog.Car, seller *sales.Seller) (*sales.Listing, error)
}

type DealerListingFactory struct{}

func (DealerListingFactory) Create(id sales.ListingID, car *catalog.Car, seller *sales.Seller) (*sales.Listing, error) {
	l, err := sales.NewListing(id, car, seller, car.Price())
	if err != nil { return nil, err }
	_ = l.SetTitle(fmt.Sprintf("%s %s %d, от дилера", car.Brand(), car.Model(), car.Year().Value()))
	_ = l.SetRegion("Moscow")
	c, _ := sales.NewDiscountPercent(3)
	_ = l.SetCommission(c)
	l.WithWarranty()
	l.WithTestDrive()
	_ = l.AddPhoto("https://example.com/placeholder.jpg")
	return l, nil
}

type PrivateListingFactory struct{}

func (PrivateListingFactory) Create(id sales.ListingID, car *catalog.Car, seller *sales.Seller) (*sales.Listing, error) {
	l, err := sales.NewListing(id, car, seller, car.Price())
	if err != nil { return nil, err }
	_ = l.SetTitle(fmt.Sprintf("%s %s %d, пробег %d", car.Brand(), car.Model(), car.Year().Value(), car.Mileage().Value()))
	_ = l.SetRegion("Moscow")
	c, _ := sales.NewDiscountPercent(1)
	_ = l.SetCommission(c)
	l.WithTestDrive()
	_ = l.AddPhoto("https://example.com/placeholder.jpg")
	return l, nil
}

type AuctionListingFactory struct{}

func (AuctionListingFactory) Create(id sales.ListingID, car *catalog.Car, seller *sales.Seller) (*sales.Listing, error) {
	l, err := sales.NewListing(id, car, seller, car.Price())
	if err != nil { return nil, err }
	_ = l.SetTitle(fmt.Sprintf("Аукцион: %s %s %d", car.Brand(), car.Model(), car.Year().Value()))
	_ = l.SetRegion("Moscow")
	c, _ := sales.NewDiscountPercent(5)
	_ = l.SetCommission(c)
	_ = l.AddPhoto("https://example.com/placeholder.jpg")
	return l, nil
}

func FactoryFor(t sales.SellerType) ListingFactory {
	switch t {
	case sales.SellerDealer:  return DealerListingFactory{}
	case sales.SellerPrivate: return PrivateListingFactory{}
	case sales.SellerAuction: return AuctionListingFactory{}
	default:
		panic("unknown seller: " + string(t))
	}
}
