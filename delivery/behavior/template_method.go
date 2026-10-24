package behavior

import (
	"log"
	"delivery/domain"
)

type Pipeline interface {
	Pickup(*domain.Delivery)
	Transport(*domain.Delivery)
	Dropoff(*domain.Delivery)
	Notify(*domain.Delivery)
}

func RunPipeline(p Pipeline, d *domain.Delivery) {
	p.Pickup(d); p.Transport(d); p.Dropoff(d); p.Notify(d)
}

type CityPipeline struct{}
func (CityPipeline) Pickup(d *domain.Delivery)    { log.Printf("[City] pickup %s", d.ID()) }
func (CityPipeline) Transport(d *domain.Delivery) { log.Printf("[City] transport %s", d.ID()) }
func (CityPipeline) Dropoff(d *domain.Delivery)   { log.Printf("[City] dropoff %s", d.ID()) }
func (CityPipeline) Notify(d *domain.Delivery)    { log.Printf("[City] SMS %s", d.ID()) }

type IntercityPipeline struct{}
func (IntercityPipeline) Pickup(d *domain.Delivery)    { log.Printf("[Inter] pickup %s", d.ID()) }
func (IntercityPipeline) Transport(d *domain.Delivery) { log.Printf("[Inter] transport %s", d.ID()) }
func (IntercityPipeline) Dropoff(d *domain.Delivery)   { log.Printf("[Inter] dropoff %s", d.ID()) }
func (IntercityPipeline) Notify(d *domain.Delivery)    { log.Printf("[Inter] push %s", d.ID()) }

type IntlPipeline struct{}
func (IntlPipeline) Pickup(d *domain.Delivery)    { log.Printf("[Intl] customs %s", d.ID()) }
func (IntlPipeline) Transport(d *domain.Delivery) { log.Printf("[Intl] transport %s", d.ID()) }
func (IntlPipeline) Dropoff(d *domain.Delivery)   { log.Printf("[Intl] dropoff %s", d.ID()) }
func (IntlPipeline) Notify(d *domain.Delivery)    { log.Printf("[Intl] email %s", d.ID()) }

func PipelineFor(t domain.DeliveryType) Pipeline {
	switch t {
	case domain.TypeIntl: return IntlPipeline{}
	case domain.TypeExpress, domain.TypeEconomy: return CityPipeline{}
	}
	return IntercityPipeline{}
}
