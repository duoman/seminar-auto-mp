package listing

import (
	"errors"

	"marketplace/domain/catalog"
	"marketplace/domain/sales"
	"marketplace/domain/shared"
)

var (
	ErrTitle = errors.New("title required")
	ErrPrice = errors.New("price must be positive")
)

type ListingBuilder struct {
	id, title, desc, region string
	car                     *catalog.Car
	seller                  *sales.Seller
	photos                  []string
	price                   shared.Money
	warranty, testDrive     bool
	commission              sales.DiscountPercent
}

func NewListingBuilder(id sales.ListingID, car *catalog.Car, seller *sales.Seller) *ListingBuilder {
	return &ListingBuilder{id: id.String(), car: car, seller: seller, price: car.Price(), region: "Moscow"}
}
func (b *ListingBuilder) Title(v string) *ListingBuilder           { b.title = v; return b }
func (b *ListingBuilder) Description(v string) *ListingBuilder     { b.desc = v; return b }
func (b *ListingBuilder) Photo(u string) *ListingBuilder           { b.photos = append(b.photos, u); return b }
func (b *ListingBuilder) Region(v string) *ListingBuilder          { b.region = v; return b }
func (b *ListingBuilder) Commission(d sales.DiscountPercent) *ListingBuilder { b.commission = d; return b }
func (b *ListingBuilder) WithWarranty() *ListingBuilder            { b.warranty = true; return b }
func (b *ListingBuilder) WithTestDrive() *ListingBuilder           { b.testDrive = true; return b }

func (b *ListingBuilder) Build() (*sales.Listing, error) {
	if b.title == "" { return nil, ErrTitle }
	if b.price.Amount() <= 0 { return nil, ErrPrice }
	id, _ := sales.NewListingID(b.id)
	l, err := sales.NewListing(id, b.car, b.seller, b.price)
	if err != nil { return nil, err }
	if err := l.SetTitle(b.title); err != nil { return nil, err }
	if b.desc != "" { _ = l.SetDescription(b.desc) }
	for _, p := range b.photos { _ = l.AddPhoto(p) }
	_ = l.SetRegion(b.region)
	if b.commission.Value() > 0 { _ = l.SetCommission(b.commission) }
	if b.warranty { l.WithWarranty() }
	if b.testDrive { l.WithTestDrive() }
	return l, nil
}
