package main

import (
	"fmt"
	"log"
	"net/http"
)

func logRequest(r *http.Request) {
	requestLine := fmt.Sprintf("%s %s %s %s %s", r.RemoteAddr, r.Host, r.Method, r.RequestURI, r.URL)
	log.Printf("%s", requestLine)
}
