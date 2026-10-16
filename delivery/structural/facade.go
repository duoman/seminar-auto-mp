package structural

import (
	"time"

	"delivery/behavior"
	"delivery/domain"
	"delivery/ports"
)

type ScheduleFacade struct {
	repo      ports.DeliveryRepository
	geo       GeoProvider
	templates *RouteTemplateFactory
	pub       *behavior.Publisher
	history   *behavior.History
}

func NewScheduleFacade(
	repo ports.DeliveryRepository,
	geo GeoProvider,
	templates *RouteTemplateFactory,
	pub *behavior.Publisher,
	history *behavior.History,
) *ScheduleFacade {
	return &ScheduleFacade{repo: repo, geo: geo, templates: templates, pub: pub, history: history}
}

func (f *ScheduleFacade) Schedule(
	id, listingID, customerID string,
	deliveryType domain.DeliveryType,
	from, to domain.Point,
) (*domain.Delivery, error) {

	strategy := behavior.StrategyFor(deliveryType)
	route := strategy.BuildRoute(from, to)
	template := f.templates.Get(from.City, to.City)

	d := &domain.Delivery{
		ID: id, ListingID: listingID, CustomerID: customerID,
		Type: deliveryType, From: from, To: to,
		Status: domain.StatusCreated,
		Price: route.PriceRUB, DistanceKM: template.DistanceKM,
		CreatedAt: time.Now(),
	}

	state := behavior.StateFor(d.Status)
	if err := state.Schedule(d); err != nil {
		return nil, err
	}
	if err := f.repo.Save(d); err != nil {
		return nil, err
	}
	f.pub.Publish(behavior.DeliveryEvent{DeliveryID: d.ID, Status: string(d.Status)})
	return d, nil
}
