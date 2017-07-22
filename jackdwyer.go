package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	_ "image/jpeg"
	_ "image/png"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var host string = "127.0.0.1"
var port string = "5000"
var indexPaginationRE = regexp.MustCompile("^/([0-9]+)$")
var deleteIdRE = regexp.MustCompile("^/delete/([0-9]+)$")

func main() {
	hostAndPort := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting Web Server @ http://%s", hostAndPort)
	r := mux.NewRouter()

	r.HandleFunc("/", index)
	r.HandleFunc("/admin", admin)
	r.HandleFunc("/delete/{id:[0-9]+}", deleteLocation)
	r.HandleFunc("/upload", upload)
	// http.HandleFunc("/delete/", deleteLocation)
	// http.HandleFunc("/upload", upload)
	// http.HandleFunc("/upload/", upload)
	// http.HandleFunc("/favicon.ico", favicon)
	db, _ = sql.Open("sqlite3", "./db/app.db")
	err := http.ListenAndServe(hostAndPort, r)
	if err != nil {
		log.Println(err)
	}
	db.Close()
}

func favicon(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	http.ServeFile(w, r, "assets/favicon.ico")
}

func admin(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	page := r.URL.Query().Get("page")
	i, err := strconv.Atoi(page)
	cur := i
	if err != nil {
		cur = 0
	}
	results, err := getLocations(cur, 30, false)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	t := template.Must(template.ParseFiles("templates/admin.tpl.html"))
	t.Execute(w, results)
}

// index/<pagination>: paginate location/image results
func index(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	var err error
	paginationValue := 0
	// this is for getting: /<pagination value>
	match := indexPaginationRE.FindStringSubmatch(r.URL.Path)
	if len(match) == 2 {
		paginationValue, err = strconv.Atoi(match[1])
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
	results, err := getLocations(paginationValue, 10, true)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	templateData := struct {
		Results  []location
		BasePath string
	}{
		results,
		basePath,
	}
	log.Printf("%T\n", results)
	t := template.Must(template.ParseFiles("templates/index.tpl.html"))
	t.Execute(w, templateData)
}

func deleteLocation(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	match := deleteIdRE.FindStringSubmatch(r.URL.Path)
	log.Println(match)
	if len(match) != 2 {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	locationId, err := strconv.Atoi(match[1])
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	log.Printf("Location id: %d\n", locationId)
	err = deleteRow(locationId)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	http.Redirect(w, r, "/admin", 301)
}

// upload: allows me to upload a new image
func upload(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	t := time.Now()
	filename := fmt.Sprintf("%s.JPG", t.Format(imageFilenameFormat))
	timestamp := fmt.Sprintf("%s", t.Format(SQLTimestampFormat))
	log.Printf("Generated filename: %s\n", filename)
	if r.Method != "POST" {
		http.Error(w, "404", 404)
		return
	}
	latitude, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "400", 400)
		return
	}
	longitude, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "400", 400)
		return
	}
	address, err := reverseGeocode(latitude, longitude)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	formFile, _, err := r.FormFile("img")
	defer formFile.Close()
	if err != nil {
		log.Printf("FAILED on image upload: %s", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	b, err := ioutil.ReadAll(formFile)
	newImage, err := ResizeImage(b)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	buff := new(bytes.Buffer)
	err = jpeg.Encode(buff, newImage, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = UploadFile(buff.Bytes(), fmt.Sprintf("960/%s", filename))
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = insertLocation(latitude, longitude, filename, address, timestamp)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	io.WriteString(w, "Success")
	return
}
