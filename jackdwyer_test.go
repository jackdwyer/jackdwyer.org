package main_test

import (
	"image"
	"os"
	"testing"

	_ "image/jpeg"

	jackdwyer "github.com/jackdwyer/jackdwyer.org"
)

func TestImageIsTooBig(t *testing.T) {
	file, _ := os.Open("./test_data/2160x1140.png")
	img, _, _ := image.DecodeConfig(file)
	if !jackdwyer.ImageTooBig(img) {
		t.Error("Image is too big")
	}
}

func TestImageIsRightSize(t *testing.T) {
	file, _ := os.Open("./test_data/960x540.png")
	img, _, _ := image.DecodeConfig(file)
	if jackdwyer.ImageTooBig(img) {
		t.Error("Image is correct width")
	}
}
