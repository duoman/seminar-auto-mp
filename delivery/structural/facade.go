package structural

import (
	"delivery/behavior"
	"delivery/domain"
	"delivery/ports"
)

type ScheduleFacade struct {
	repo ports.Repository
	tpl  *Templates
	pub  *behavior.Publisher
}

func NewScheduleFacade(repo ports.Repository, tpl *Templates, pub *behavior.Publisher) *ScheduleFacade {
	return &ScheduleFacade{repo: repo, tpl: tpl, pub: pub}
}

func (f *ScheduleFacade) Schedule(id, listingID, customerID, typ, from, to string) (*domain.Delivery, error) {
	strategy := behavior.StrategyFor(typ)
	_, price, distance := strategy.BuildRoute(from, to)
	t := f.tpl.Get(from, to)
	if t.DistanceKM > 0 { distance = t.DistanceKM }

	d, err := domain.NewDelivery(id, listingID, customerID, domain.DeliveryType(typ),
		domain.Point{City: from}, domain.Point{City: to}, price, distance)
	if err != nil { return nil, err }
	if err := d.Schedule(); err != nil { return nil, err }
	if err := f.repo.Save(d); err != nil { return nil, err }
	f.pub.Publish(behavior.Event{DeliveryID: d.ID(), Status: string(d.Status())})
	return d, nil
}
