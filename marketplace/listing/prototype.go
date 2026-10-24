package listing

import "marketplace/domain/sales"

func Clone(src *sales.Listing) (*sales.Listing, error) {
	if src == nil { return nil, nil }
	newID, _ := sales.NewListingID(src.ID().String() + "-copy")
	c, err := sales.NewListing(newID, src.Car(), src.Seller(), src.Price())
	if err != nil { return nil, err }
	if src.Title() != "" { _ = c.SetTitle(src.Title() + " (копия)") }
	if src.Description() != "" { _ = c.SetDescription(src.Description()) }
	for _, p := range src.Photos() { _ = c.AddPhoto(p) }
	_ = c.SetRegion(src.Region())
	_ = c.SetCommission(src.Commission())
	if src.HasWarranty() { c.WithWarranty() }
	if src.HasTestDrive() { c.WithTestDrive() }
	return c, nil
}
