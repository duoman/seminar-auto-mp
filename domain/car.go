package domain

type Car struct {
	ID      string
	Brand   string
	Model   string
	Year    int
	Mileage int
	Price   int
}

func NewCar(id, brand, model string, year, mileage, price int) *Car {
	return &Car{ID: id, Brand: brand, Model: model, Year: year, Mileage: mileage, Price: price}
}
