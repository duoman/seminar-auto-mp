package domain

import "fmt"

type DeliveryType string
const (
	TypeExpress DeliveryType = "express"
	TypeEconomy DeliveryType = "economy"
	TypeEco     DeliveryType = "eco"
	TypeIntl    DeliveryType = "international"
)

type DeliveryStatus string
const (
	StatusCreated   DeliveryStatus = "created"
	StatusScheduled DeliveryStatus = "scheduled"
	StatusInTransit DeliveryStatus = "in_transit"
	StatusDelivered DeliveryStatus = "delivered"
	StatusCancelled DeliveryStatus = "cancelled"
)

type Point struct{ City, Address string }

type Delivery struct {
	id, listingID, customerID string
	typ                        DeliveryType
	from, to                   Point
	status                     DeliveryStatus
	price, distanceKM          int
}

func NewDelivery(id, listingID, customerID string, typ DeliveryType, from, to Point, price, distance int) (*Delivery, error) {
	if id == "" || listingID == "" { return nil, fmt.Errorf("ids required") }
	switch typ {
	case TypeExpress, TypeEconomy, TypeEco, TypeIntl:
	default:
		return nil, fmt.Errorf("unknown type: %s", typ)
	}
	return &Delivery{id, listingID, customerID, typ, from, to, StatusCreated, price, distance}, nil
}

func Restore(id, listingID, customerID string, typ DeliveryType, from, to Point, status DeliveryStatus, price, distance int) *Delivery {
	return &Delivery{id, listingID, customerID, typ, from, to, status, price, distance}
}

func (d *Delivery) ID() string             { return d.id }
func (d *Delivery) ListingID() string      { return d.listingID }
func (d *Delivery) CustomerID() string     { return d.customerID }
func (d *Delivery) Type() DeliveryType     { return d.typ }
func (d *Delivery) From() Point            { return d.from }
func (d *Delivery) To() Point              { return d.to }
func (d *Delivery) Status() DeliveryStatus { return d.status }
func (d *Delivery) Price() int             { return d.price }
func (d *Delivery) DistanceKM() int        { return d.distanceKM }

func (d *Delivery) Schedule() error {
	if d.status != StatusCreated { return fmt.Errorf("only from created, now: %s", d.status) }
	d.status = StatusScheduled
	return nil
}
func (d *Delivery) Start() error {
	if d.status != StatusScheduled { return fmt.Errorf("only from scheduled, now: %s", d.status) }
	d.status = StatusInTransit
	return nil
}
func (d *Delivery) Complete() error {
	if d.status != StatusInTransit { return fmt.Errorf("only from in_transit, now: %s", d.status) }
	d.status = StatusDelivered
	return nil
}
