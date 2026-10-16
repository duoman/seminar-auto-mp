package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPDeliveryClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPDeliveryClient(baseURL string) *HTTPDeliveryClient {
	return &HTTPDeliveryClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

type deliveryRequest struct {
	ListingID  string `json:"listing_id"`
	CustomerID string `json:"customer_id"`
	Type       string `json:"type"`
	FromCity   string `json:"from_city"`
	ToCity     string `json:"to_city"`
}

type deliveryResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (c *HTTPDeliveryClient) Schedule(listingID, customerID, fromCity, toCity, deliveryType string) (string, error) {
	body, _ := json.Marshal(deliveryRequest{
		ListingID: listingID, CustomerID: customerID,
		Type: deliveryType, FromCity: fromCity, ToCity: toCity,
	})
	resp, err := c.client.Post(c.baseURL+"/api/v1/deliveries", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("delivery call: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("delivery returned %d", resp.StatusCode)
	}
	var dr deliveryResponse
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		return "", err
	}
	return dr.ID, nil
}
