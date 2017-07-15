package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

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
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/delete", deleteLocation)
	http.HandleFunc("/delete/", deleteLocation)
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/upload/", upload)
	http.HandleFunc("/", index)
	db, _ = sql.Open("sqlite3", "./db/app.db")
	http.ListenAndServe(hostAndPort, nil)
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
	t := template.Must(template.ParseFiles("templates/index.tpl.html"))
	t.Execute(w, results)
}

func deleteLocation(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	match := deleteIdRE.FindStringSubmatch(r.URL.Path)
	fmt.Println(match)
	if len(match) != 2 {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	locationId, err := strconv.Atoi(match[1])
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	fmt.Printf("Location id: %d\n", locationId)
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
	if r.Method != "POST" {
		http.Error(w, "404", 404)
		return
	}
	file, handler, err := r.FormFile("img")
	defer file.Close()
	if err != nil {
		log.Printf("FAILED on image upload: %s", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0444)
	defer f.Close()
	if err != nil {
		log.Printf("FAILED in saving image: %s", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	io.Copy(f, file)
}
