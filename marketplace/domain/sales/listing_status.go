package sales

type ListingStatus string
const (
	StatusDraft    ListingStatus = "draft"
	StatusPending  ListingStatus = "pending"
	StatusActive   ListingStatus = "active"
	StatusRejected ListingStatus = "rejected"
	StatusSold     ListingStatus = "sold"
)
