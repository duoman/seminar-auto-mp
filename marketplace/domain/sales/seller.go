package sales

import "fmt"

type SellerID struct{ value string }
func NewSellerID(v string) (SellerID, error) {
	if v == "" { return SellerID{}, fmt.Errorf("seller id required") }
	return SellerID{value: v}, nil
}
func (id SellerID) String() string { return id.value }

type SellerType string
const (
	SellerDealer  SellerType = "dealer"
	SellerPrivate SellerType = "private"
	SellerAuction SellerType = "auction"
)

type Seller struct {
	id   SellerID
	name string
	typ  SellerType
}

func NewSeller(id SellerID, name string, typ SellerType) (*Seller, error) {
	if name == "" { return nil, fmt.Errorf("name required") }
	switch typ {
	case SellerDealer, SellerPrivate, SellerAuction:
	default:
		return nil, fmt.Errorf("unknown type: %s", typ)
	}
	return &Seller{id, name, typ}, nil
}
func (s *Seller) ID() SellerID    { return s.id }
func (s *Seller) Name() string     { return s.name }
func (s *Seller) Type() SellerType { return s.typ }
