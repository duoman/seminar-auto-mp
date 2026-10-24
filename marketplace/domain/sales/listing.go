package sales

import (
	"fmt"

	"marketplace/domain/catalog"
	"marketplace/domain/shared"
)

type Listing struct {
	id           ListingID
	car          *catalog.Car
	seller       *Seller
	title        string
	description  string
	photos       []string
	price        shared.Money
	status       ListingStatus
	hasWarranty  bool
	hasTestDrive bool
	region       string
	commission   DiscountPercent
	deliveryID   string
}

func NewListing(id ListingID, car *catalog.Car, seller *Seller, price shared.Money) (*Listing, error) {
	if car == nil || seller == nil { return nil, fmt.Errorf("car and seller required") }
	return &Listing{id: id, car: car, seller: seller, price: price, status: StatusDraft}, nil
}

// Restore — используется репозиторием при восстановлении из БД.
func Restore(id ListingID, car *catalog.Car, seller *Seller, title, description string,
	photos []string, price shared.Money, status ListingStatus,
	warranty, testDrive bool, region string, commission DiscountPercent, deliveryID string) *Listing {
	return &Listing{
		id: id, car: car, seller: seller, title: title, description: description,
		photos: photos, price: price, status: status,
		hasWarranty: warranty, hasTestDrive: testDrive, region: region,
		commission: commission, deliveryID: deliveryID,
	}
}

func (l *Listing) ID() ListingID              { return l.id }
func (l *Listing) Car() *catalog.Car          { return l.car }
func (l *Listing) Seller() *Seller            { return l.seller }
func (l *Listing) Title() string              { return l.title }
func (l *Listing) Description() string        { return l.description }
func (l *Listing) Photos() []string           { return l.photos }
func (l *Listing) Price() shared.Money        { return l.price }
func (l *Listing) Status() ListingStatus      { return l.status }
func (l *Listing) HasWarranty() bool          { return l.hasWarranty }
func (l *Listing) HasTestDrive() bool         { return l.hasTestDrive }
func (l *Listing) Region() string             { return l.region }
func (l *Listing) Commission() DiscountPercent { return l.commission }
func (l *Listing) DeliveryID() string         { return l.deliveryID }

func (l *Listing) SetTitle(v string) error {
	if l.status != StatusDraft { return fmt.Errorf("edit only in draft, current: %s", l.status) }
	if len(v) < 3 { return fmt.Errorf("title too short") }
	l.title = v
	return nil
}
func (l *Listing) SetDescription(v string) error {
	if l.status != StatusDraft { return fmt.Errorf("edit only in draft") }
	l.description = v
	return nil
}
func (l *Listing) AddPhoto(url string) error {
	if l.status != StatusDraft { return fmt.Errorf("edit only in draft") }
	if len(l.photos) >= 10 { return fmt.Errorf("max 10 photos") }
	l.photos = append(l.photos, url)
	return nil
}
func (l *Listing) SetRegion(v string) error {
	if l.status != StatusDraft { return fmt.Errorf("edit only in draft") }
	l.region = v
	return nil
}
func (l *Listing) SetCommission(d DiscountPercent) error {
	if l.status != StatusDraft { return fmt.Errorf("edit only in draft") }
	l.commission = d
	return nil
}
func (l *Listing) WithWarranty()  { l.hasWarranty = true }
func (l *Listing) WithTestDrive() { l.hasTestDrive = true }

func (l *Listing) Submit() error {
	if l.status != StatusDraft { return fmt.Errorf("only from draft, current: %s", l.status) }
	if l.title == "" { return fmt.Errorf("title required") }
	l.status = StatusPending
	return nil
}
func (l *Listing) Publish() error {
	if l.status != StatusPending { return fmt.Errorf("only from pending, current: %s", l.status) }
	if l.title == "" { return fmt.Errorf("title required") }
	if len(l.photos) == 0 { return fmt.Errorf("photo required") }
	l.status = StatusActive
	return nil
}
func (l *Listing) MarkSold() error {
	if l.status != StatusActive { return fmt.Errorf("only active can sell, current: %s", l.status) }
	l.status = StatusSold
	return nil
}
func (l *Listing) AttachDelivery(deliveryID string) error {
	if l.status != StatusSold { return fmt.Errorf("delivery only for sold, current: %s", l.status) }
	if deliveryID == "" { return fmt.Errorf("delivery id required") }
	l.deliveryID = deliveryID
	return nil
}
