package main

import (
	"image"

	"github.com/nfnt/resize"
)

func ImageTooBig(img image.Config) bool {
	if img.Width > resizeWidth {
		return true
	}
	return false
}

func ResizeImage(img image.Image) (image.Image, error) {
	return resize.Resize(uint(resizeWidth), 0, img, resize.Lanczos3), nil
}
