package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nfnt/resize"
)

func ImageTooBig(img image.Config) bool {
	if img.Width > resizeWidth {
		return true
	}
	return false
}

func logRequest(r *http.Request) {
	requestLine := fmt.Sprintf("%s %s %s %s %s", r.RemoteAddr, r.Host, r.Method, r.RequestURI, r.URL)
	log.Printf("%s", requestLine)
}

func ResizeImage(img image.Image) (image.Image, error) {
	return resize.Resize(uint(resizeWidth), 0, img, resize.Lanczos3), nil
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
		Bucket:        aws.String("dev.jackdwyer.org"),
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
