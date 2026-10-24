package shared

import "fmt"

type Money struct {
	amount   int64
	currency string
}

func NewMoney(kopecks int64, currency string) (Money, error) {
	if kopecks < 0 {
		return Money{}, fmt.Errorf("amount must be non-negative")
	}
	if currency == "" {
		return Money{}, fmt.Errorf("currency is required")
	}
	return Money{amount: kopecks, currency: currency}, nil
}
func (m Money) Amount() int64    { return m.amount }
func (m Money) Currency() string { return m.currency }
func (m Money) String() string   { return fmt.Sprintf("%d.%02d %s", m.amount/100, m.amount%100, m.currency) }
