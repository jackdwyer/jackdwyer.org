package main

import (
	"image"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"

	_ "image/jpeg"
	_ "image/png"
)

func HandleUpload(imgPath string) {
	file, _ := os.Open(imgPath)
	defer file.Close()
	img, _, _ := image.Decode(file)
	m := resize.Resize(960, 0, img, resize.Lanczos3)
	out, _ := os.Create("/tmp/out.jpg")
	defer out.Close()
	jpeg.Encode(out, m, nil)
}
