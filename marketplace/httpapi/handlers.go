package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"marketplace/domain/catalog"
	"marketplace/domain/sales"
	"marketplace/domain/shared"
	"marketplace/listing"
	"marketplace/ports"
)

type Server struct {
	repo ports.ListingRepository
}

func NewServer(repo ports.ListingRepository) *Server { return &Server{repo: repo} }

func writeJSON(w http.ResponseWriter, c int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, c int, m string) {
	writeJSON(w, c, map[string]string{"error": m})
}

type createReq struct {
	CarID      string `json:"car_id"`
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	Year       int    `json:"year"`
	Mileage    int    `json:"mileage"`
	Price      int    `json:"price"`
	SellerID   string `json:"seller_id"`
	SellerName string `json:"seller_name"`
	SellerType string `json:"seller_type"`
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeErr(w, 400, "invalid json"); return }

	carID, err := catalog.NewCarID(req.CarID); if err != nil { writeErr(w, 400, err.Error()); return }
	year, err := catalog.NewYear(req.Year); if err != nil { writeErr(w, 400, err.Error()); return }
	mileage, err := catalog.NewMileage(req.Mileage); if err != nil { writeErr(w, 400, err.Error()); return }
	price, err := shared.NewMoney(int64(req.Price)*100, "RUB"); if err != nil { writeErr(w, 400, err.Error()); return }
	car, err := catalog.NewCar(carID, req.Brand, req.Model, year, mileage, price); if err != nil { writeErr(w, 400, err.Error()); return }
	sID, _ := sales.NewSellerID(req.SellerID)
	seller, err := sales.NewSeller(sID, req.SellerName, sales.SellerType(req.SellerType))
	if err != nil { writeErr(w, 400, err.Error()); return }

	lID, _ := sales.NewListingID("L-" + req.SellerID + "-" + req.CarID)
	l, err := listing.FactoryFor(seller.Type()).Create(lID, car, seller)
	if err != nil { writeErr(w, 400, err.Error()); return }
	if err := l.Submit(); err != nil { writeErr(w, 400, err.Error()); return }

	if err := s.repo.Save(l); err != nil { writeErr(w, 500, err.Error()); return }
	writeJSON(w, 201, toResp(l))
}

func (s *Server) Publish(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/listings/"), "/publish")
	lID, _ := sales.NewListingID(id)
	l, err := s.repo.FindByID(lID)
	if err != nil { writeErr(w, 404, "not found"); return }
	if err := l.Publish(); err != nil { writeErr(w, 400, err.Error()); return }
	_ = s.repo.Save(l)
	writeJSON(w, 200, toResp(l))
}

// Sell — публикует событие ListingSold в outbox.
// Никаких HTTP-вызовов к delivery. Только локальная транзакция + событие.
func (s *Server) Sell(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/listings/"), "/sell")
	lID, _ := sales.NewListingID(id)
	l, err := s.repo.FindByID(lID)
	if err != nil { writeErr(w, 404, "not found"); return }
	if err := l.MarkSold(); err != nil { writeErr(w, 400, err.Error()); return }

	payload, _ := json.Marshal(map[string]any{
		"event_id":    fmt.Sprintf("evt-listing-sold-%s-%d", l.ID().String(), time.Now().UnixNano()),
		"listing_id":  l.ID().String(),
		"customer_id": "CUST-1",
		"from_city":   "Moscow",
		"to_city":     "SPb",
		"type":        "express",
	})
	evt := sales.OutboxEvent{
		ID:      fmt.Sprintf("outbox-%s-%d", l.ID().String(), time.Now().UnixNano()),
		Type:    sales.EventListingSold,
		Payload: payload,
	}

	if err := s.repo.SaveWithOutbox(l, evt); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, toResp(l))
}

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/listings/")
	lID, _ := sales.NewListingID(id)
	l, err := s.repo.FindByID(lID)
	if err != nil { writeErr(w, 404, "not found"); return }
	writeJSON(w, 200, toResp(l))
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	all, err := s.repo.All()
	if err != nil { writeErr(w, 500, err.Error()); return }
	out := make([]any, 0, len(all))
	for _, l := range all { out = append(out, toResp(l)) }
	writeJSON(w, 200, out)
}

func toResp(l *sales.Listing) map[string]any {
	return map[string]any{
		"id":             l.ID().String(),
		"title":          l.Title(),
		"description":    l.Description(),
		"photos":         l.Photos(),
		"price":          l.Price().Amount() / 100,
		"status":         string(l.Status()),
		"commission":     l.Commission().Value(),
		"has_warranty":   l.HasWarranty(),
		"has_test_drive": l.HasTestDrive(),
		"region":         l.Region(),
		"delivery_id":    l.DeliveryID(),
	}
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/listings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet: s.List(w, r)
		case http.MethodPost: s.Create(w, r)
		}
	})
	mux.HandleFunc("/api/v1/listings/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/publish"): s.Publish(w, r)
		case strings.HasSuffix(r.URL.Path, "/sell"):    s.Sell(w, r)
		case r.Method == http.MethodGet:                s.Get(w, r)
		}
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})
	return mux
}
