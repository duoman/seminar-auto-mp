package behavior

import (
	"log"
	"delivery/domain"
)

type Kind interface{ Accept(Visitor) }
type LocalKind struct{ D *domain.Delivery }
func (k LocalKind) Accept(v Visitor) { v.VisitLocal(k) }
type LongHaulKind struct{ D *domain.Delivery }
func (k LongHaulKind) Accept(v Visitor) { v.VisitLongHaul(k) }
type IntlKind struct{ D *domain.Delivery }
func (k IntlKind) Accept(v Visitor) { v.VisitIntl(k) }

type Visitor interface {
	VisitLocal(LocalKind)
	VisitLongHaul(LongHaulKind)
	VisitIntl(IntlKind)
}

type InsuranceVisitor struct{ Total int }
func (v *InsuranceVisitor) VisitLocal(k LocalKind) { add := k.D.Price() / 100; v.Total += add; log.Printf("[Insurance] local %s +%d", k.D.ID(), add) }
func (v *InsuranceVisitor) VisitLongHaul(k LongHaulKind) { add := k.D.Price() / 50; v.Total += add; log.Printf("[Insurance] long %s +%d", k.D.ID(), add) }
func (v *InsuranceVisitor) VisitIntl(k IntlKind) { add := k.D.Price() / 20; v.Total += add; log.Printf("[Insurance] intl %s +%d", k.D.ID(), add) }

type ReportVisitor struct{ Lines []string }
func (v *ReportVisitor) VisitLocal(k LocalKind) { v.Lines = append(v.Lines, k.D.ID()) }
func (v *ReportVisitor) VisitLongHaul(k LongHaulKind) { v.Lines = append(v.Lines, k.D.ID()) }
func (v *ReportVisitor) VisitIntl(k IntlKind) { v.Lines = append(v.Lines, k.D.ID()) }

func Classify(d *domain.Delivery) Kind {
	if d.Type() == domain.TypeIntl { return IntlKind{d} }
	if d.From().City == d.To().City { return LocalKind{d} }
	return LongHaulKind{d}
}
