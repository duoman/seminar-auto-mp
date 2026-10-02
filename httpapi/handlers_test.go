package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"marketplace/httpapi"
	"marketplace/infra"
)

func newTestServer() http.Handler {
	repo := infra.NewInMemoryListingRepository()
	return httpapi.NewServer(repo).Routes()
}

func TestCreateListing_Dealer(t *testing.T) {
	srv := newTestServer()

	body := `{
		"car_id": "CAR-1", "brand": "Tesla", "model": "Model 3",
		"year": 2023, "mileage": 15000, "price": 5000000,
		"seller_id": "D-1", "seller_name": "AutoDealer", "seller_type": "dealer"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/listings", strings.NewReader(body))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)

	if resp["commission"].(float64) != 3 {
		t.Fatalf("want commission 3, got %v", resp["commission"])
	}
	if !strings.Contains(resp["title"].(string), "официального дилера") {
		t.Fatalf("unexpected title: %s", resp["title"])
	}
}

func TestCloneListing_CopiesPhotos(t *testing.T) {
	srv := newTestServer()

	create := `{
		"car_id": "CAR-1", "brand": "Tesla", "model": "Model 3",
		"year": 2023, "mileage": 15000, "price": 5000000,
		"seller_id": "D-1", "seller_name": "AutoDealer", "seller_type": "dealer",
		"title": "Tesla", "description": "desc", "region": "Moscow",
		"photos": ["https://example.com/1.jpg"],
		"has_warranty": true, "has_test_drive": true, "commission": 3
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/listings/custom", strings.NewReader(create))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create failed: %d %s", rec.Code, rec.Body.String())
	}

	var created map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	id := created["id"].(string)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/listings/"+id+"/clone", nil)
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("clone failed: %d %s", rec.Code, rec.Body.String())
	}

	var clone map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &clone)

	photos := clone["photos"].([]any)
	if len(photos) != 1 {
		t.Fatalf("clone must keep 1 photo, got %d", len(photos))
	}
}

func TestPublishListing_Activates(t *testing.T) {
	srv := newTestServer()

	create := `{
		"car_id": "CAR-1", "brand": "Tesla", "model": "Model 3",
		"year": 2023, "mileage": 15000, "price": 5000000,
		"seller_id": "D-1", "seller_name": "AutoDealer", "seller_type": "dealer"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/listings", strings.NewReader(create))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	var created map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	id := created["id"].(string)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/listings/"+id+"/publish", nil)
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("publish failed: %d %s", rec.Code, rec.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["status"].(string) != "active" {
		t.Fatalf("want active, got %s", resp["status"])
	}
}

func TestMarketplaces_ReturnsThreeFamilies(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/marketplaces", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	var items []map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &items)
	if len(items) != 3 {
		t.Fatalf("want 3 marketplaces, got %d", len(items))
	}
}
