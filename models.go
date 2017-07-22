package main

import (
	"database/sql"
	"log"
)

type location struct {
	Id           int
	uuid         string
	_type        string
	Latitude     float64
	Longitude    float64
	accuracy     float32
	Timestamp    string
	Image        sql.NullString
	comment      []byte
	address      string
	ShortAddress sql.NullString
	Unlocked     bool
}

func deleteRow(id int) error {
	_, err := db.Exec("delete from location where id = ?", id)
	return err
}

func insertLocation(latitude float64, longitude float64, image string, address string, timestamp string) error {
	statment := "INSERT INTO location (latitude, longitude, image, short_address, timestamp, unlocked) VALUES (?, ?, ?, ?, ?, ?);"
	res, err := db.Exec(statment, latitude, longitude, image, address, timestamp, true)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	log.Printf("Inserted Location with ID: %d", id)
	return nil
}

func getLocations(offsetBase int, limit int, unlocked bool) ([]location, error) {
	var result location
	var results []location
	if unlocked {
		offset := offsetBase * 10
		rows, err := db.Query("select image, timestamp, short_address from location where unlocked=1 order by id desc limit ? offset ?;", limit, offset)
		defer rows.Close()
		if err != nil {
			log.Printf("Failed Query: %q", err)
		}
		for rows.Next() {
			rows.Scan(&result.Image, &result.Timestamp, &result.ShortAddress)
			results = append(results, result)
		}
		return results, err
	}
	offset := offsetBase * 30
	rows, err := db.Query("select id, latitude, longitude, image, timestamp, short_address, unlocked from location order by id desc limit ? offset ?;", limit, offset)
	defer rows.Close()
	if err != nil {
		log.Printf("Failed Query: %q", err)
	}
	for rows.Next() {
		rows.Scan(&result.Id, &result.Latitude, &result.Longitude, &result.Image, &result.Timestamp, &result.ShortAddress, &result.Unlocked)
		results = append(results, result)
	}
	return results, err
}
