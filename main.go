package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Reservation struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	Guests int    `json:"guests"`
	Phone  string `json:"phone"`
}

func createReservation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var res Reservation
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	sqlStatement :=
		`
		INSERT INTO reservations (name, date, time, guests, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`
	id := 0
	err := db.QueryRow(sqlStatement, res.Name, res.Date, res.Time, res.Guests, res.Phone).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID int `json:"id"`
	}{ID: id})
}

func getReservations(w http.ResponseWriter, r *http.Request) {
	//
}

func getReservation(w http.ResponseWriter, r *http.Request) {
	//
}

func updateReservation(w http.ResponseWriter, r *http.Request) {
	//
}

func deleteReservation(w http.ResponseWriter, r *http.Request) {
	//
}

func main() {
	// Database connection
	dbConfig := DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "myuser",
		Password: "mypassword",
		DBName:   "dinner_reservations",
	}
	db := NewDB(dbConfig)
	defer db.Close()

	// HTTP Handlers
	http.HandleFunc("/reservations/create", func(w http.ResponseWriter, r *http.Request) {
		createReservation(db, w, r)
	})

	// Start the HTTP server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
