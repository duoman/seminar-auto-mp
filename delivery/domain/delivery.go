package domain

import "time"

type DeliveryType string

const (
	TypeExpress       DeliveryType = "express"
	TypeEconomy       DeliveryType = "economy"
	TypeEco           DeliveryType = "eco"
	TypeInternational DeliveryType = "international"
)

type DeliveryStatus string

const (
	StatusCreated   DeliveryStatus = "created"
	StatusScheduled DeliveryStatus = "scheduled"
	StatusInTransit DeliveryStatus = "in_transit"
	StatusDelivered DeliveryStatus = "delivered"
	StatusCancelled DeliveryStatus = "cancelled"
)

type Point struct {
	City    string
	Address string
}

type Route struct {
	Description string
	DistanceKM  int
	PriceRUB    int
}

type Delivery struct {
	ID         string
	ListingID  string
	CustomerID string
	Type       DeliveryType
	From       Point
	To         Point
	Status     DeliveryStatus
	Price      int
	DistanceKM int
	CreatedAt  time.Time
}
