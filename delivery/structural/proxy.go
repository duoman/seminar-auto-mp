package structural

import "sync"

type LazyGeoProxy struct {
	once sync.Once
	real GeoProvider
	sdk  MapsSDK
}

func NewLazyGeoProxy(s MapsSDK) *LazyGeoProxy { return &LazyGeoProxy{sdk: s} }
func (p *LazyGeoProxy) DistanceKM(from, to string) int {
	p.once.Do(func() { p.real = NewMapsAdapter(p.sdk) })
	return p.real.DistanceKM(from, to)
}
