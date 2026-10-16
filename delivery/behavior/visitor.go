package behavior

import (
	"fmt"

	"delivery/domain"
)

type DeliveryKind interface{ Accept(v DeliveryVisitor) }

type LocalKind struct{ Delivery *domain.Delivery }
func (k LocalKind) Accept(v DeliveryVisitor) { v.VisitLocal(k) }

type LongHaulKind struct{ Delivery *domain.Delivery }
func (k LongHaulKind) Accept(v DeliveryVisitor) { v.VisitLongHaul(k) }

type InternationalKind struct{ Delivery *domain.Delivery }
func (k InternationalKind) Accept(v DeliveryVisitor) { v.VisitInternational(k) }

type DeliveryVisitor interface {
	VisitLocal(LocalKind)
	VisitLongHaul(LongHaulKind)
	VisitInternational(InternationalKind)
}

type InsuranceVisitor struct{ Total int }

func (v *InsuranceVisitor) VisitLocal(k LocalKind) {
	add := k.Delivery.Price / 100
	v.Total += add
	fmt.Printf("  [Insurance] local %s: +%d\n", k.Delivery.ID, add)
}
func (v *InsuranceVisitor) VisitLongHaul(k LongHaulKind) {
	add := k.Delivery.Price / 50
	v.Total += add
	fmt.Printf("  [Insurance] long-haul %s: +%d\n", k.Delivery.ID, add)
}
func (v *InsuranceVisitor) VisitInternational(k InternationalKind) {
	add := k.Delivery.Price / 20
	v.Total += add
	fmt.Printf("  [Insurance] international %s: +%d\n", k.Delivery.ID, add)
}

type CostReportVisitor struct{ Lines []string }

func (v *CostReportVisitor) VisitLocal(k LocalKind) {
	v.Lines = append(v.Lines, fmt.Sprintf("Local %s: %d", k.Delivery.ID, k.Delivery.Price))
}
func (v *CostReportVisitor) VisitLongHaul(k LongHaulKind) {
	v.Lines = append(v.Lines, fmt.Sprintf("Long-haul %s: %d", k.Delivery.ID, k.Delivery.Price))
}
func (v *CostReportVisitor) VisitInternational(k InternationalKind) {
	v.Lines = append(v.Lines, fmt.Sprintf("International %s: %d", k.Delivery.ID, k.Delivery.Price))
}

func ClassifyDelivery(d *domain.Delivery) DeliveryKind {
	switch {
	case d.Type == domain.TypeInternational:
		return InternationalKind{Delivery: d}
	case d.From.City == d.To.City:
		return LocalKind{Delivery: d}
	default:
		return LongHaulKind{Delivery: d}
	}
}
