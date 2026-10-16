package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"marketplace/domain"
	"marketplace/listing"
	"marketplace/ports"
)

type Server struct {
	repo     ports.ListingRepository
	delivery ports.DeliveryClient
}

func NewServer(repo ports.ListingRepository, delivery ports.DeliveryClient) *Server {
	return &Server{repo: repo, delivery: delivery}
}

// ---------- DTO ----------

type createListingRequest struct {
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

type customListingRequest struct {
	createListingRequest
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Region       string   `json:"region"`
	Photos       []string `json:"photos"`
	HasWarranty  bool     `json:"has_warranty"`
	HasTestDrive bool     `json:"has_test_drive"`
	Commission   int      `json:"commission"`
}

type listingResponse struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Photos       []string `json:"photos"`
	Price        int      `json:"price"`
	Status       string   `json:"status"`
	Commission   int      `json:"commission"`
	HasWarranty  bool     `json:"has_warranty"`
	HasTestDrive bool     `json:"has_test_drive"`
	Region       string   `json:"region"`
	DeliveryID   string   `json:"delivery_id,omitempty"`
}

func toResponse(l *domain.Listing) listingResponse {
	return listingResponse{
		ID:           l.ID,
		Title:        l.Title,
		Description:  l.Description,
		Photos:       l.Photos,
		Price:        l.Price,
		Status:       string(l.Status),
		Commission:   l.Commission,
		HasWarranty:  l.HasWarranty,
		HasTestDrive: l.HasTestDrive,
		Region:       l.Region,
		DeliveryID:   l.DeliveryID,
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

// ---------- Handlers ----------

// POST /api/v1/listings — Factory Method
func (s *Server) CreateListing(w http.ResponseWriter, r *http.Request) {
	var req createListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.SellerType == "" {
		writeError(w, http.StatusBadRequest, "seller_type is required")
		return
	}

	car := domain.NewCar(req.CarID, req.Brand, req.Model, req.Year, req.Mileage, req.Price)
	seller := &domain.Seller{
		ID:   req.SellerID,
		Name: req.SellerName,
		Type: domain.SellerType(req.SellerType),
	}

	factory, err := listing.FactoryFor(seller.Type)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	l := factory.Create(car, seller)

	if err := s.repo.Save(l); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(l))
}

// POST /api/v1/listings/custom — Builder
func (s *Server) CreateCustomListing(w http.ResponseWriter, r *http.Request) {
	var req customListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	car := domain.NewCar(req.CarID, req.Brand, req.Model, req.Year, req.Mileage, req.Price)
	seller := &domain.Seller{
		ID:   req.SellerID,
		Name: req.SellerName,
		Type: domain.SellerType(req.SellerType),
	}

	b := listing.NewListingBuilder(car, seller).
		Title(req.Title).
		Description(req.Description).
		Region(req.Region).
		Commission(req.Commission)

	for _, p := range req.Photos {
		b.Photo(p)
	}
	if req.HasWarranty {
		b.WithWarranty()
	}
	if req.HasTestDrive {
		b.WithTestDrive()
	}

	l, err := b.Build()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.repo.Save(l); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(l))
}

// POST /api/v1/listings/{id}/clone — Prototype
func (s *Server) CloneListing(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/listings/"), "/clone")

	original, err := s.repo.FindByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "listing not found")
		return
	}
	clone := listing.Clone(original)
	if err := s.repo.Save(clone); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(clone))
}

// POST /api/v1/listings/{id}/publish
func (s *Server) PublishListing(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/listings/"), "/publish")

	l, err := s.repo.FindByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "listing not found")
		return
	}
	if strings.TrimSpace(l.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	if l.Price <= 0 {
		writeError(w, http.StatusBadRequest, "price must be positive")
		return
	}
	l.Status = domain.StatusActive
	if err := s.repo.Save(l); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toResponse(l))
}

// POST /api/v1/listings/{id}/sell — DIP + DI: вызывает delivery через интерфейс
func (s *Server) SellListing(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/listings/"), "/sell")

	l, err := s.repo.FindByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "listing not found")
		return
	}
	if l.Status != domain.StatusActive {
		writeError(w, http.StatusBadRequest, "listing must be active")
		return
	}

	// Вызов delivery-микросервиса через интерфейс DeliveryClient
	deliveryID, err := s.delivery.Schedule(l.ID, "CUST-1", "Moscow", "SPb", "express")
	if err != nil {
		writeError(w, http.StatusBadGateway, "delivery failed: "+err.Error())
		return
	}

	l.Status = domain.StatusSold
	l.DeliveryID = deliveryID
	if err := s.repo.Save(l); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toResponse(l))
}

// GET /api/v1/listings/{id}
func (s *Server) GetListing(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/listings/")
	l, err := s.repo.FindByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "listing not found")
		return
	}
	writeJSON(w, http.StatusOK, toResponse(l))
}

// GET /api/v1/listings
func (s *Server) ListListings(w http.ResponseWriter, r *http.Request) {
	all, err := s.repo.All()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]listingResponse, 0, len(all))
	for _, l := range all {
		out = append(out, toResponse(l))
	}
	writeJSON(w, http.StatusOK, out)
}

// GET /api/v1/marketplaces — Abstract Factory
func (s *Server) ListMarketplaces(w http.ResponseWriter, r *http.Request) {
	factories := map[string]listing.MarketplaceFactory{
		"standard": listing.StandardFactory{},
		"premium":  listing.PremiumFactory{},
		"auction":  listing.AuctionFactory{},
	}
	type item struct {
		Name         string `json:"name"`
		Commission   int    `json:"commission"`
		ManualReview bool   `json:"manual_review"`
	}
	out := make([]item, 0, len(factories))
	for name, f := range factories {
		out = append(out, item{
			Name:         name,
			Commission:   f.CreateCommission().Percent(),
			ManualReview: f.CreateModeration().RequiresManualReview(),
		})
	}
	writeJSON(w, http.StatusOK, out)
}

// ---------- Router ----------

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/listings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.ListListings(w, r)
		case http.MethodPost:
			s.CreateListing(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/api/v1/listings/custom", s.CreateCustomListing)

	mux.HandleFunc("/api/v1/listings/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/clone") && r.Method == http.MethodPost:
			s.CloneListing(w, r)
		case strings.HasSuffix(r.URL.Path, "/publish") && r.Method == http.MethodPost:
			s.PublishListing(w, r)
		case strings.HasSuffix(r.URL.Path, "/sell") && r.Method == http.MethodPost:
			s.SellListing(w, r)
		case r.Method == http.MethodGet:
			s.GetListing(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/api/v1/marketplaces", s.ListMarketplaces)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return mux
}
