package behavior

import (
	"fmt"

	"delivery/domain"
)

type DeliveryState interface {
	Name() string
	Schedule(d *domain.Delivery) error
	Start(d *domain.Delivery) error
	Complete(d *domain.Delivery) error
	Cancel(d *domain.Delivery) error
}

type CreatedState struct{}
func (CreatedState) Name() string { return "created" }
func (CreatedState) Schedule(d *domain.Delivery) error { d.Status = domain.StatusScheduled; return nil }
func (s CreatedState) Start(*domain.Delivery) error    { return fmt.Errorf("сначала запланировать (%s)", s.Name()) }
func (s CreatedState) Complete(*domain.Delivery) error { return fmt.Errorf("ещё не в пути (%s)", s.Name()) }
func (CreatedState) Cancel(d *domain.Delivery) error   { d.Status = domain.StatusCancelled; return nil }

type ScheduledState struct{}
func (ScheduledState) Name() string { return "scheduled" }
func (ScheduledState) Schedule(*domain.Delivery) error { return fmt.Errorf("уже запланировано") }
func (ScheduledState) Start(d *domain.Delivery) error  { d.Status = domain.StatusInTransit; return nil }
func (s ScheduledState) Complete(*domain.Delivery) error { return fmt.Errorf("ещё не в пути (%s)", s.Name()) }
func (ScheduledState) Cancel(d *domain.Delivery) error { d.Status = domain.StatusCancelled; return nil }

type InTransitState struct{}
func (InTransitState) Name() string { return "in_transit" }
func (s InTransitState) Schedule(*domain.Delivery) error { return fmt.Errorf("уже в пути (%s)", s.Name()) }
func (s InTransitState) Start(*domain.Delivery) error    { return fmt.Errorf("уже в пути (%s)", s.Name()) }
func (InTransitState) Complete(d *domain.Delivery) error { d.Status = domain.StatusDelivered; return nil }
func (s InTransitState) Cancel(*domain.Delivery) error   { return fmt.Errorf("нельзя отменить (%s)", s.Name()) }

type DeliveredState struct{}
func (DeliveredState) Name() string { return "delivered" }
func (s DeliveredState) Schedule(*domain.Delivery) error { return fmt.Errorf("уже доставлено (%s)", s.Name()) }
func (s DeliveredState) Start(*domain.Delivery) error    { return fmt.Errorf("уже доставлено (%s)", s.Name()) }
func (s DeliveredState) Complete(*domain.Delivery) error { return fmt.Errorf("уже доставлено (%s)", s.Name()) }
func (s DeliveredState) Cancel(*domain.Delivery) error   { return fmt.Errorf("уже доставлено (%s)", s.Name()) }

func StateFor(s domain.DeliveryStatus) DeliveryState {
	switch s {
	case domain.StatusCreated:
		return CreatedState{}
	case domain.StatusScheduled:
		return ScheduledState{}
	case domain.StatusInTransit:
		return InTransitState{}
	case domain.StatusDelivered:
		return DeliveredState{}
	default:
		return nil
	}
}
