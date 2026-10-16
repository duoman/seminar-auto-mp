package structural

import (
	"fmt"
	"time"
)

type RouteCalculator interface{ Calc(from, to string) int }

type SimpleCalculator struct{}
func (SimpleCalculator) Calc(from, to string) int { return 700 }

type LoggingDecorator struct{ wrappee RouteCalculator }

func NewLoggingDecorator(w RouteCalculator) LoggingDecorator { return LoggingDecorator{wrappee: w} }

func (d LoggingDecorator) Calc(from, to string) int {
	start := time.Now()
	r := d.wrappee.Calc(from, to)
	fmt.Printf("  [Logging] %s→%s = %d in %v\n", from, to, r, time.Since(start))
	return r
}

type CachingDecorator struct {
	wrappee RouteCalculator
	cache   map[string]int
}

func NewCachingDecorator(w RouteCalculator) *CachingDecorator {
	return &CachingDecorator{wrappee: w, cache: map[string]int{}}
}

func (d *CachingDecorator) Calc(from, to string) int {
	key := from + "|" + to
	if v, ok := d.cache[key]; ok {
		fmt.Printf("  [Caching] hit %s\n", key)
		return v
	}
	v := d.wrappee.Calc(from, to)
	d.cache[key] = v
	return v
}
