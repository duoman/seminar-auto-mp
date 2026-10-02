package listing

import "marketplace/domain"

func Clone(src *domain.Listing) *domain.Listing {
	if src == nil {
		return nil
	}
	clone := *src
	if src.Car != nil {
		c := *src.Car
		clone.Car = &c
	}
	if src.Seller != nil {
		s := *src.Seller
		clone.Seller = &s
	}
	if src.Photos != nil {
		clone.Photos = make([]string, len(src.Photos))
		copy(clone.Photos, src.Photos)
	}
	clone.ID = src.ID + "-copy"
	clone.Title = src.Title + " (копия)"
	clone.Status = domain.StatusDraft
	return &clone
}
