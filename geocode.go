package main

import (
	_ "github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/openstreetmap"
)

func reverseGeocode(latitude float64, longitude float64) (string, error) {
	geo := openstreetmap.Geocoder()
	address, err := geo.ReverseGeocode(latitude, longitude)
	if err != nil {
		return "", err
	}
	return address.FormattedAddress, nil
}
