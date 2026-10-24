package domain

// Анемичная модель: публичные поля, никаких методов, никакой защиты.
// Всё, что могло бы быть инвариантом, — «в голове разработчика».
type Listing struct {
	ID           string
	CarID        string
	Title        string
	Photos       []string
	Price        int    // ← можно -1000
	Status       string // ← можно "несуществующий"
	Commission   int    // ← можно 150
	HasWarranty  bool
	HasTestDrive bool
	Region       string
}
