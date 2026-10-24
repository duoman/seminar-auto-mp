package catalog

import (
	"fmt"

	"marketplace/domain/shared"
)

type CarID struct{ value string }

func NewCarID(v string) (CarID, error) {
	if v == "" { return CarID{}, fmt.Errorf("car id required") }
	return CarID{value: v}, nil
}
func (c CarID) String() string { return c.value }

type Year struct{ value int }
func NewYear(v int) (Year, error) {
	if v < 1900 || v > 2100 { return Year{}, fmt.Errorf("year out of range: %d", v) }
	return Year{value: v}, nil
}
func (y Year) Value() int { return y.value }

type Mileage struct{ value int }
func NewMileage(v int) (Mileage, error) {
	if v < 0 { return Mileage{}, fmt.Errorf("mileage negative") }
	return Mileage{value: v}, nil
}
func (m Mileage) Value() int { return m.value }

type Car struct {
	id      CarID
	brand   string
	model   string
	year    Year
	mileage Mileage
	price   shared.Money
}

func NewCar(id CarID, brand, model string, year Year, mileage Mileage, price shared.Money) (*Car, error) {
	if brand == "" || model == "" { return nil, fmt.Errorf("brand and model required") }
	return &Car{id, brand, model, year, mileage, price}, nil
}
func (c *Car) ID() CarID           { return c.id }
func (c *Car) Brand() string       { return c.brand }
func (c *Car) Model() string       { return c.model }
func (c *Car) Year() Year          { return c.year }
func (c *Car) Mileage() Mileage    { return c.mileage }
func (c *Car) Price() shared.Money { return c.price }
