package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

func ImageTooBig(img image.Config) bool {
	if img.Width > resizeWidth {
		return true
	}
	return false
}

func ResizeImage(b []byte) (image.Image, error) {
	r := bytes.NewReader(b)
	image, err := jpeg.Decode(r)
	if err != nil {
		return nil, nil
	}
	_, err = r.Seek(0, 0)
	if err != nil {
		return nil, nil
	}
	x, _ := exif.Decode(r)
	orientation, _ := x.Get(exif.Orientation)
	o := orientation.String()
	log.Printf("Orientation of image: %s\n", o)
	if o == "6" {
		nimg := imaging.Resize(image, 0, resizeWidth, imaging.Lanczos)
		return imaging.Rotate270(nimg), nil
	}
	return imaging.Resize(image, resizeWidth, 0, imaging.Lanczos), nil
}

func UploadFile(f []byte, filename string) error {
	creds := credentials.NewEnvCredentials()
	_, err := creds.Get()
	if err != nil {
		return err
	}
	cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds)
	awsS3 := s3.New(session.New(), cfg)
	size := len(f)
	fileType := http.DetectContentType(f)
	params := &s3.PutObjectInput{
		ACL:           &S3ACL,
		Bucket:        aws.String("dev-images.jackdwyer.org"),
		Key:           aws.String(filename),
		Body:          bytes.NewReader(f),
		ContentLength: aws.Int64(int64(size)),
		ContentType:   aws.String(fileType),
	}
	_, err = awsS3.PutObject(params)
	if err != nil {
		return err
	}
	return nil
}
