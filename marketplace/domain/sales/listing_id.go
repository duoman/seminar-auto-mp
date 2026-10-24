package sales

import "fmt"

type ListingID struct{ value string }
func NewListingID(v string) (ListingID, error) {
	if v == "" { return ListingID{}, fmt.Errorf("listing id required") }
	return ListingID{value: v}, nil
}
func (id ListingID) String() string { return id.value }
