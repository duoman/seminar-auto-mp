package structural

import "log"

type Calculator interface{ Calc(from, to string) int }
type SimpleCalc struct{}
func (SimpleCalc) Calc(_, _ string) int { return 700 }

type LogDecorator struct{ W Calculator }
func NewLog(w Calculator) LogDecorator { return LogDecorator{w} }
func (d LogDecorator) Calc(f, t string) int {
	r := d.W.Calc(f, t)
	log.Printf("[Log] %s→%s = %d", f, t, r)
	return r
}

type CacheDecorator struct {
	W Calculator
	M map[string]int
}
func NewCache(w Calculator) *CacheDecorator { return &CacheDecorator{W: w, M: map[string]int{}} }
func (d *CacheDecorator) Calc(f, t string) int {
	k := f + "|" + t
	if v, ok := d.M[k]; ok { return v }
	v := d.W.Calc(f, t)
	d.M[k] = v
	return v
}
