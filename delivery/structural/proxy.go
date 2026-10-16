package structural

import (
	"fmt"
	"sync"
)

type LazyGeoProxy struct {
	once sync.Once
	real GeoProvider
	sdk  MapsSDK
}

func NewLazyGeoProxy(sdk MapsSDK) *LazyGeoProxy { return &LazyGeoProxy{sdk: sdk} }

func (p *LazyGeoProxy) DistanceKM(from, to string) int {
	p.once.Do(func() {
		fmt.Println("  [LazyGeoProxy] first call — creating adapter")
		p.real = NewMapsAdapter(p.sdk)
	})
	return p.real.DistanceKM(from, to)
}
