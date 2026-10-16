package structural

import "fmt"

type MapsSDK struct{}

func (MapsSDK) CalcRoute(fromLat, fromLng, toLat, toLng float64) map[string]any {
	return map[string]any{"distance_m": 700000, "duration_s": 21600}
}

type GeoProvider interface {
	DistanceKM(from, to string) int
}

type MapsAdapter struct{ sdk MapsSDK }

func NewMapsAdapter(sdk MapsSDK) MapsAdapter { return MapsAdapter{sdk: sdk} }

func (a MapsAdapter) DistanceKM(from, to string) int {
	raw := a.sdk.CalcRoute(0, 0, 0, 0)
	meters, _ := raw["distance_m"].(int)
	fmt.Printf("  [MapsAdapter] %s→%s = %d km\n", from, to, meters/1000)
	return meters / 1000
}
