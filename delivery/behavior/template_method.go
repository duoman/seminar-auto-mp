package behavior

import (
	"fmt"

	"delivery/domain"
)

type DeliveryPipeline interface {
	Pickup(d *domain.Delivery)
	Transport(d *domain.Delivery)
	Dropoff(d *domain.Delivery)
	Notify(d *domain.Delivery)
}

func RunDelivery(p DeliveryPipeline, d *domain.Delivery) {
	p.Pickup(d); p.Transport(d); p.Dropoff(d); p.Notify(d)
}

type CityPipeline struct{}
func (CityPipeline) Pickup(d *domain.Delivery)    { fmt.Printf("  [City] курьер забрал %s\n", d.ID) }
func (CityPipeline) Transport(d *domain.Delivery) { fmt.Printf("  [City] доставка внутри %s\n", d.To.City) }
func (CityPipeline) Dropoff(d *domain.Delivery)   { fmt.Printf("  [City] курьер отдал %s\n", d.ID) }
func (CityPipeline) Notify(d *domain.Delivery)    { fmt.Printf("  [City] SMS по %s\n", d.ID) }

type IntercityPipeline struct{}
func (IntercityPipeline) Pickup(d *domain.Delivery)    { fmt.Printf("  [Intercity] склад принял %s\n", d.ID) }
func (IntercityPipeline) Transport(d *domain.Delivery) { fmt.Printf("  [Intercity] перевозка %s→%s\n", d.From.City, d.To.City) }
func (IntercityPipeline) Dropoff(d *domain.Delivery)   { fmt.Printf("  [Intercity] выдача %s\n", d.ID) }
func (IntercityPipeline) Notify(d *domain.Delivery)    { fmt.Printf("  [Intercity] push %s\n", d.ID) }

type InternationalPipeline struct{}
func (InternationalPipeline) Pickup(d *domain.Delivery)    { fmt.Printf("  [Intl] таможня приняла %s\n", d.ID) }
func (InternationalPipeline) Transport(d *domain.Delivery) { fmt.Printf("  [Intl] перевозка %s→%s\n", d.From.City, d.To.City) }
func (InternationalPipeline) Dropoff(d *domain.Delivery)   { fmt.Printf("  [Intl] выдача %s\n", d.ID) }
func (InternationalPipeline) Notify(d *domain.Delivery)    { fmt.Printf("  [Intl] email %s\n", d.ID) }

func PipelineFor(t domain.DeliveryType) DeliveryPipeline {
	switch t {
	case domain.TypeInternational:
		return InternationalPipeline{}
	case domain.TypeExpress, domain.TypeEconomy:
		return CityPipeline{}
	default:
		return IntercityPipeline{}
	}
}
