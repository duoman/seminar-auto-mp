package sales

import "fmt"

type DiscountPercent struct{ value int }
func NewDiscountPercent(v int) (DiscountPercent, error) {
	if v < 0 || v > 100 { return DiscountPercent{}, fmt.Errorf("discount out of range: %d", v) }
	return DiscountPercent{value: v}, nil
}
func (d DiscountPercent) Value() int    { return d.value }
func (d DiscountPercent) String() string { return fmt.Sprintf("%d%%", d.value) }
