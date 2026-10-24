package behavior

import "fmt"

type RouteStrategy interface{ BuildRoute(from, to string) (string, int, int) }

type FastestRoute struct{}
func (FastestRoute) BuildRoute(f, t string) (string, int, int) { return fmt.Sprintf("Fast %s→%s", f, t), 15000, 700 }
type CheapestRoute struct{}
func (CheapestRoute) BuildRoute(f, t string) (string, int, int) { return fmt.Sprintf("Cheap %s→%s", f, t), 8000, 800 }
type EcoRoute struct{}
func (EcoRoute) BuildRoute(f, t string) (string, int, int) { return fmt.Sprintf("Eco %s→%s", f, t), 12000, 750 }
type IntlRoute struct{}
func (IntlRoute) BuildRoute(f, t string) (string, int, int) { return fmt.Sprintf("Intl %s→%s", f, t), 60000, 2500 }

func StrategyFor(t string) RouteStrategy {
	switch t {
	case "express": return FastestRoute{}
	case "economy": return CheapestRoute{}
	case "eco": return EcoRoute{}
	case "international": return IntlRoute{}
	}
	return FastestRoute{}
}
