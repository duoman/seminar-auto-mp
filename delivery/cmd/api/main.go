package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"delivery/behavior"
	"delivery/httpapi"
	"delivery/infra"
	"delivery/structural"
)

func main() {
	repo := infra.NewInMemoryDeliveryRepository()

	pub := &behavior.Publisher{}
	pub.Subscribe(behavior.CustomerObserver{})
	pub.Subscribe(behavior.WarehouseObserver{})
	pub.Subscribe(behavior.AnalyticsObserver{})
	pub.Subscribe(behavior.MarketplaceObserver{})

	history := &behavior.History{}
	geo := structural.NewLazyGeoProxy(structural.MapsSDK{})
	templates := structural.NewRouteTemplateFactory()
	facade := structural.NewScheduleFacade(repo, geo, templates, pub, history)

	var calc structural.RouteCalculator = structural.SimpleCalculator{}
	calc = structural.NewLoggingDecorator(calc)
	cached := structural.NewCachingDecorator(calc)
	fmt.Println("=== Decorator demo ===")
	_ = cached.Calc("Moscow", "SPb")
	_ = cached.Calc("Moscow", "SPb")

	server := httpapi.NewServer(repo, facade, pub, history)
	addr := ":" + envOr("PORT", "8081")
	log.Printf("delivery on %s", addr)
	if err := http.ListenAndServe(addr, server.Routes()); err != nil {
		log.Fatal(err)
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
