package behavior

import "fmt"

type DeliveryEvent struct {
	DeliveryID string
	Status     string
}

type Observer interface{ OnDeliveryChanged(e DeliveryEvent) }

type Publisher struct{ observers []Observer }

func (p *Publisher) Subscribe(o Observer) { p.observers = append(p.observers, o) }
func (p *Publisher) Publish(e DeliveryEvent) {
	for _, o := range p.observers {
		o.OnDeliveryChanged(e)
	}
}

type CustomerObserver struct{}
func (CustomerObserver) OnDeliveryChanged(e DeliveryEvent) {
	fmt.Printf("  [Customer] %s → %s\n", e.DeliveryID, e.Status)
}

type WarehouseObserver struct{}
func (WarehouseObserver) OnDeliveryChanged(e DeliveryEvent) {
	fmt.Printf("  [Warehouse] %s → %s\n", e.DeliveryID, e.Status)
}

type AnalyticsObserver struct{}
func (AnalyticsObserver) OnDeliveryChanged(e DeliveryEvent) {
	fmt.Printf("  [Analytics] %s → %s\n", e.DeliveryID, e.Status)
}

type MarketplaceObserver struct{}
func (MarketplaceObserver) OnDeliveryChanged(e DeliveryEvent) {
	fmt.Printf("  [Marketplace] %s → %s\n", e.DeliveryID, e.Status)
}
