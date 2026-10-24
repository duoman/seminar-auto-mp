package structural

type MapsSDK struct{}
func (MapsSDK) CalcRoute(_, _, _, _ float64) map[string]any { return map[string]any{"distance_m": 700000} }

type GeoProvider interface{ DistanceKM(from, to string) int }

type MapsAdapter struct{ sdk MapsSDK }
func NewMapsAdapter(s MapsSDK) MapsAdapter { return MapsAdapter{sdk: s} }
func (a MapsAdapter) DistanceKM(from, to string) int {
	r := a.sdk.CalcRoute(0, 0, 0, 0)
	m, _ := r["distance_m"].(int)
	return m / 1000
}
