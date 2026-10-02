package domain

type SellerType string

const (
	SellerDealer  SellerType = "dealer"
	SellerPrivate SellerType = "private"
	SellerAuction SellerType = "auction"
)

type Seller struct {
	ID     string
	Name   string
	Type   SellerType
	Rating float64
}
