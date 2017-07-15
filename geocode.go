package main

import (
	"fmt"

	_ "github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/openstreetmap"
)

func testGeocode() {
	lat, lng := -37.813611, 144.963056
	geo := openstreetmap.Geocoder()
	address, _ := geo.ReverseGeocode(lat, lng)
	fmt.Println(address)
}
