package domain

type ListingStatus string

const (
	StatusDraft    ListingStatus = "draft"
	StatusPending  ListingStatus = "pending"
	StatusActive   ListingStatus = "active"
	StatusRejected ListingStatus = "rejected"
	StatusSold     ListingStatus = "sold"
)

type Listing struct {
	ID           string
	Car          *Car
	Seller       *Seller
	Title        string
	Description  string
	Photos       []string
	Price        int
	Status       ListingStatus
	HasWarranty  bool
	HasTestDrive bool
	Region       string
	Commission   int
	DeliveryID   string
}
