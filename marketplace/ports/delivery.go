package ports

type DeliveryClient interface {
	Schedule(listingID, customerID, fromCity, toCity, deliveryType string) (string, error)
}
