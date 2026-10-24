package behavior

import "log"

type Event struct{ DeliveryID, Status string }
type Observer interface{ Notify(Event) }

type Publisher struct{ obs []Observer }
func (p *Publisher) Subscribe(o Observer) { p.obs = append(p.obs, o) }
func (p *Publisher) Publish(e Event) {
	for _, o := range p.obs { o.Notify(e) }
}

type LogObserver struct{ Name string }
func (o LogObserver) Notify(e Event) {
	log.Printf("[%s] delivery %s → %s", o.Name, e.DeliveryID, e.Status)
}
