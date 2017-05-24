package main_test

import (
	"bufio"
	"bytes"
	"image"
	"io/ioutil"
	"os"
	"testing"

	jackdwyer "github.com/jackdwyer/jackdwyer.org"
	jpeg "image/jpeg"
	// png "image/png"
)

func fail(t *testing.T, err error) {
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestImageIsTooBig(t *testing.T) {
	file, err := os.Open("./test_data/2160x1440.png")
	fail(t, err)
	img, _, err := image.DecodeConfig(file)
	fail(t, err)
	if !jackdwyer.ImageTooBig(img) {
		t.Error("Image is too big")
	}
}

func TestImageIsRightSize(t *testing.T) {
	file, err := os.Open("./test_data/960x540.png")
	fail(t, err)
	img, _, err := image.DecodeConfig(file)
	fail(t, err)
	if jackdwyer.ImageTooBig(img) {
		t.Error("Image is correct width")
	}
}

func TestImageResize(t *testing.T) {
	var buff bytes.Buffer
	writer := bufio.NewWriter(&buff)
	file, err := os.Open("./test_data/2160x1440.png")
	defer file.Close()
	fail(t, err)
	img, _, err := image.Decode(file)
	fail(t, err)
	resizedImage, _ := jackdwyer.ResizeImage(img)
	err = jpeg.Encode(writer, resizedImage, nil)
	fail(t, err)
	resizedImageConfig, _, err := image.DecodeConfig(bufio.NewReader(&buff))
	fail(t, err)
	if jackdwyer.ImageTooBig(resizedImageConfig) {
		t.Error("Image was not resized")
	}
}

func TestFileUpload(t *testing.T) {
	fileByteArr, err := ioutil.ReadFile("./test_data/2160x1440.png")
	fail(t, err)
	_ = jackdwyer.UploadFile(fileByteArr, "thisissometestfile.png")
}
