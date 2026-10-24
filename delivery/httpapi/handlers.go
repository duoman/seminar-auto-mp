package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"delivery/behavior"
	"delivery/domain"
	"delivery/ports"
	"delivery/structural"
)

type Server struct {
	repo   ports.Repository
	facade *structural.ScheduleFacade
	pub    *behavior.Publisher
}

func NewServer(r ports.Repository, f *structural.ScheduleFacade, p *behavior.Publisher) *Server {
	return &Server{repo: r, facade: f, pub: p}
}

func writeJSON(w http.ResponseWriter, c int, v any) {
	w.Header().Set("Content-Type", "application/json"); w.WriteHeader(c); _ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, c int, m string) { writeJSON(w, c, map[string]string{"error": m}) }

type dResp struct {
	ID         string `json:"id"`
	ListingID  string `json:"listing_id"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
	Type       string `json:"type"`
	DistanceKM int    `json:"distance_km"`
}

func toResp(d *domain.Delivery) dResp {
	return dResp{d.ID(), d.ListingID(), string(d.Status()), d.Price(), string(d.Type()), d.DistanceKM()}
}

func (s *Server) Start(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/deliveries/"), "/start")
	d, err := s.repo.FindByID(id)
	if err != nil { writeErr(w, 404, "not found"); return }
	if err := d.Start(); err != nil { writeErr(w, 400, err.Error()); return }
	_ = s.repo.Save(d)
	behavior.RunPipeline(behavior.PipelineFor(d.Type()), d)
	s.pub.Publish(behavior.Event{DeliveryID: d.ID(), Status: string(d.Status())})
	writeJSON(w, 200, toResp(d))
}

func (s *Server) Complete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/deliveries/"), "/complete")
	d, err := s.repo.FindByID(id)
	if err != nil { writeErr(w, 404, "not found"); return }
	if err := d.Complete(); err != nil { writeErr(w, 400, err.Error()); return }
	_ = s.repo.Save(d)
	s.pub.Publish(behavior.Event{DeliveryID: d.ID(), Status: string(d.Status())})
	writeJSON(w, 200, toResp(d))
}

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/deliveries/")
	d, err := s.repo.FindByID(id)
	if err != nil { writeErr(w, 404, "not found"); return }
	writeJSON(w, 200, toResp(d))
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	all, _ := s.repo.All()
	out := make([]dResp, 0, len(all))
	for _, d := range all { out = append(out, toResp(d)) }
	writeJSON(w, 200, out)
}

func (s *Server) Reports(w http.ResponseWriter, r *http.Request) {
	all, _ := s.repo.All()
	ins := &behavior.InsuranceVisitor{}
	rep := &behavior.ReportVisitor{}
	for _, d := range all {
		behavior.Classify(d).Accept(ins)
		behavior.Classify(d).Accept(rep)
	}
	writeJSON(w, 200, map[string]any{"total_insurance": ins.Total, "report": rep.Lines})
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/deliveries", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet { s.List(w, r) }
	})
	mux.HandleFunc("/api/v1/deliveries/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/start"): s.Start(w, r)
		case strings.HasSuffix(r.URL.Path, "/complete"): s.Complete(w, r)
		case r.Method == http.MethodGet: s.Get(w, r)
		}
	})
	mux.HandleFunc("/api/v1/reports", s.Reports)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})
	return mux
}
