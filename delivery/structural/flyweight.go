package structural

import "log"

type RouteTemplate struct {
	FromCity, ToCity string
	DistanceKM       int
}

type Templates struct{ M map[string]*RouteTemplate }
func NewTemplates() *Templates { return &Templates{M: map[string]*RouteTemplate{}} }
func (f *Templates) Get(from, to string) *RouteTemplate {
	k := from + "|" + to
	if t, ok := f.M[k]; ok { return t }
	t := &RouteTemplate{from, to, 700}
	f.M[k] = t
	log.Printf("[Flyweight] created %s", k)
	return t
}
