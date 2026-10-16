package structural

import "fmt"

type RouteTemplate struct {
	FromCity, ToCity string
	DistanceKM       int
	BasePrice        int
}

type RouteTemplateFactory struct{ templates map[string]*RouteTemplate }

func NewRouteTemplateFactory() *RouteTemplateFactory {
	return &RouteTemplateFactory{templates: map[string]*RouteTemplate{}}
}

func (f *RouteTemplateFactory) Get(from, to string) *RouteTemplate {
	key := from + "|" + to
	if t, ok := f.templates[key]; ok {
		return t
	}
	t := &RouteTemplate{FromCity: from, ToCity: to, DistanceKM: 700, BasePrice: 10000}
	f.templates[key] = t
	fmt.Printf("  [Flyweight] created template %s\n", key)
	return t
}

func (f *RouteTemplateFactory) Count() int { return len(f.templates) }
