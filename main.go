package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
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
		log.Println("Error during reservation creation: ", err)
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID int `json:"id"`
	}{ID: id})
}

func getReservations(db *sql.DB, w http.ResponseWriter) {
	var reservations []Reservation
	sqlStatement :=
		`
		SELECT id, name, date, time, guests, phone
		FROM reservations
		`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		http.Error(w, "Failed to get reservations", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var res Reservation
		if err := rows.Scan(&res.ID, &res.Name, &res.Date, &res.Time, &res.Guests, &res.Phone); err != nil {
			http.Error(w, "Failed to get reservations", http.StatusInternalServerError)
			return
		}
		reservations = append(reservations, res)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to get reservations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservations)
}

func getReservation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}
	id := parts[2]

	var res Reservation
	sqlStatement :=
		`
    SELECT id, name, date, time, guests, phone
    FROM reservations
    WHERE id = $1
    `
	err := db.QueryRow(sqlStatement, id).Scan(&res.ID, &res.Name, &res.Date, &res.Time, &res.Guests, &res.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No reservation found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get reservation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func updateReservation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}
	id := parts[2]

	// Decode the update data from the request payload
	var updateData Reservation
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Update
	sqlStatement :=
		`
	UPDATE reservations
	SET name = $2, date = $3, time = $4, guests = $5, phone = $6
	WHERE id = $1
	`
	_, err := db.Exec(sqlStatement, id, updateData.Name, updateData.Date, updateData.Time, updateData.Guests, updateData.Phone)
	if err != nil {
		http.Error(w, "Failed to update reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{Message: "Reservation updated successfully"})
}

func deleteReservation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extract ID
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}
	id := parts[2]

	// Delete
	sqlStatement :=
		`
	DELETE FROM reservations
	WHERE id = $1
	`
	result, err := db.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, "Failed to delete reservation", http.StatusInternalServerError)
		return
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to delete reservation", http.StatusInternalServerError)
		return
	}
	if rowsDeleted == 0 {
		http.Error(w, "No reservation found with given ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{Message: "Reservation deleted successfully"})
}

func main() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection
	dbConfig := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     func() int { port, _ := strconv.Atoi(os.Getenv("DB_PORT")); return port }(),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
	db := NewDB(dbConfig)
	defer db.Close()

	// HTTP Handlers
	http.HandleFunc("/reservations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createReservation(db, w, r)
		case http.MethodGet:
			getReservations(db, w)
		default:
			http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/reservations/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getReservation(db, w, r)
		case http.MethodPut:
			updateReservation(db, w, r)
		case http.MethodDelete:
			deleteReservation(db, w, r)
		default:
			http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
