package behavior

import (
	"fmt"

	"delivery/domain"
)

type RouteStrategy interface {
	Name() string
	BuildRoute(from, to domain.Point) domain.Route
}

type FastestRoute struct{}
func (FastestRoute) Name() string { return "fastest" }
func (FastestRoute) BuildRoute(from, to domain.Point) domain.Route {
	return domain.Route{
		Description: fmt.Sprintf("Быстрый %s→%s", from.City, to.City),
		DistanceKM:  700, PriceRUB: 15000,
	}
}

type CheapestRoute struct{}
func (CheapestRoute) Name() string { return "cheapest" }
func (CheapestRoute) BuildRoute(from, to domain.Point) domain.Route {
	return domain.Route{
		Description: fmt.Sprintf("Дешёвый %s→%s", from.City, to.City),
		DistanceKM:  800, PriceRUB: 8000,
	}
}

type EcoRoute struct{}
func (EcoRoute) Name() string { return "eco" }
func (EcoRoute) BuildRoute(from, to domain.Point) domain.Route {
	return domain.Route{
		Description: fmt.Sprintf("Эко %s→%s", from.City, to.City),
		DistanceKM:  750, PriceRUB: 12000,
	}
}

type InternationalRoute struct{}
func (InternationalRoute) Name() string { return "international" }
func (InternationalRoute) BuildRoute(from, to domain.Point) domain.Route {
	return domain.Route{
		Description: fmt.Sprintf("Международный %s→%s", from.City, to.City),
		DistanceKM:  2500, PriceRUB: 60000,
	}
}

func StrategyFor(t domain.DeliveryType) RouteStrategy {
	switch t {
	case domain.TypeExpress:
		return FastestRoute{}
	case domain.TypeEconomy:
		return CheapestRoute{}
	case domain.TypeEco:
		return EcoRoute{}
	case domain.TypeInternational:
		return InternationalRoute{}
	default:
		return FastestRoute{}
	}
}
