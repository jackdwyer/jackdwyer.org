package main

import (
	"image"
)

var maxWidth = 960

func ImageTooBig(img image.Config) bool {
	if img.Width >= maxWidth {
		return false
	}
	return true
}
