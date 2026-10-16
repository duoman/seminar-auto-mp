package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"delivery/behavior"
	"delivery/domain"
	"delivery/ports"
	"delivery/structural"
)

type Server struct {
	repo    ports.DeliveryRepository
	facade  *structural.ScheduleFacade
	pub     *behavior.Publisher
	history *behavior.History
}

func NewServer(repo ports.DeliveryRepository, facade *structural.ScheduleFacade,
	pub *behavior.Publisher, history *behavior.History) *Server {
	return &Server{repo: repo, facade: facade, pub: pub, history: history}
}

type createRequest struct {
	ListingID  string `json:"listing_id"`
	CustomerID string `json:"customer_id"`
	Type       string `json:"type"`
	FromCity   string `json:"from_city"`
	ToCity     string `json:"to_city"`
}

type deliveryResponse struct {
	ID         string `json:"id"`
	ListingID  string `json:"listing_id"`
	CustomerID string `json:"customer_id"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
	FromCity   string `json:"from_city"`
	ToCity     string `json:"to_city"`
}

func toResponse(d *domain.Delivery) deliveryResponse {
	return deliveryResponse{
		ID: d.ID, ListingID: d.ListingID, CustomerID: d.CustomerID,
		Type: string(d.Type), Status: string(d.Status), Price: d.Price,
		FromCity: d.From.City, ToCity: d.To.City,
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

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid json")
		return
	}
	id := "DEL-" + req.ListingID
	from := domain.Point{City: req.FromCity}
	to := domain.Point{City: req.ToCity}
	d, err := s.facade.Schedule(id, req.ListingID, req.CustomerID, domain.DeliveryType(req.Type), from, to)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	writeJSON(w, 201, toResponse(d))
}

func (s *Server) Start(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/start")
	d, err := s.repo.FindByID(id)
	if err != nil {
		writeError(w, 404, "not found")
		return
	}
	state := behavior.StateFor(d.Status)
	if state == nil {
		writeError(w, 400, "unknown state")
		return
	}
	if err := state.Start(d); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	_ = s.repo.Save(d)
	pipeline := behavior.PipelineFor(d.Type)
	fmt.Printf("[Pipeline] %s type=%s\n", d.ID, d.Type)
	behavior.RunDelivery(pipeline, d)
	s.pub.Publish(behavior.DeliveryEvent{DeliveryID: d.ID, Status: string(d.Status)})
	writeJSON(w, 200, toResponse(d))
}

func (s *Server) Complete(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/complete")
	d, err := s.repo.FindByID(id)
	if err != nil {
		writeError(w, 404, "not found")
		return
	}
	state := behavior.StateFor(d.Status)
	if err := state.Complete(d); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	_ = s.repo.Save(d)
	s.pub.Publish(behavior.DeliveryEvent{DeliveryID: d.ID, Status: string(d.Status)})
	writeJSON(w, 200, toResponse(d))
}

func (s *Server) Undo(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/undo")
	d, err := s.repo.FindByID(id)
	if err != nil {
		writeError(w, 404, "not found")
		return
	}
	cmd := s.history.UndoLast()
	if cmd == nil {
		writeError(w, 400, "nothing to undo")
		return
	}
	_ = s.repo.Save(d)
	writeJSON(w, 200, toResponse(d))
}

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/deliveries/")
	d, err := s.repo.FindByID(id)
	if err != nil {
		writeError(w, 404, "not found")
		return
	}
	writeJSON(w, 200, toResponse(d))
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	all, _ := s.repo.All()
	out := make([]deliveryResponse, 0, len(all))
	for _, d := range all {
		out = append(out, toResponse(d))
	}
	writeJSON(w, 200, out)
}

func (s *Server) Reports(w http.ResponseWriter, r *http.Request) {
	all, _ := s.repo.All()
	insurance := &behavior.InsuranceVisitor{}
	report := &behavior.CostReportVisitor{}
	for _, d := range all {
		kind := behavior.ClassifyDelivery(d)
		kind.Accept(insurance)
		kind.Accept(report)
	}
	writeJSON(w, 200, map[string]any{
		"total_insurance": insurance.Total,
		"report":          report.Lines,
	})
}

func extractID(path, suffix string) string {
	id := strings.TrimPrefix(path, "/api/v1/deliveries/")
	return strings.TrimSuffix(id, suffix)
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/deliveries", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.List(w, r)
		case http.MethodPost:
			s.Create(w, r)
		default:
			writeError(w, 405, "method not allowed")
		}
	})
	mux.HandleFunc("/api/v1/deliveries/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/start") && r.Method == http.MethodPost:
			s.Start(w, r)
		case strings.HasSuffix(r.URL.Path, "/complete") && r.Method == http.MethodPost:
			s.Complete(w, r)
		case strings.HasSuffix(r.URL.Path, "/undo") && r.Method == http.MethodPost:
			s.Undo(w, r)
		case r.Method == http.MethodGet:
			s.Get(w, r)
		default:
			writeError(w, 405, "method not allowed")
		}
	})
	mux.HandleFunc("/api/v1/reports", s.Reports)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})
	return mux
}
