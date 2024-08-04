package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Reservation struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	Guests int    `json:"guests"`
	Phone  string `json:"phone"`
}

var (
	reservations = []Reservation{}
	nextID       = 1
	mutex        = sync.Mutex
)

func createReservation(w http.ResponseWriter, r *http.Request) {
	var reservation Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	reservation.ID = nextID
	nextID++
	reservations = append(reservations, reservation)
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reservation)
}

func getReservations(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	json.NewEncoder(w).Encode(reservations)
}

func getReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for _, reservation := range reservations {
		if reservation.ID == id {
			json.NewEncoder(w).Encode(reservation)
			return
		}
	}

	http.Error(w, "Reservation not found", http.StatusNotFound)
}

func updateReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedReservation Reservation
	if err := json.NewDecoder(r.Body).Decode(&updatedReservation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, reservation := range reservations {
		if reservation.ID == id {
			updatedReservation.ID = reservation.ID
			reservations[i] = updatedReservation
			json.NewEncoder(w).Encode(updatedReservation)
			return
		}
	}

	http.Error(w, "Reservation not found", http.StatusNotFound)
}

func deleteReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, reservation := range reservations {
		if reservation.ID == id {
			reservations = append(reservations[:i], reservations[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Reservation not found", http.StatusNotFound)
}

func main() {
	// Define routes and handlers here
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Dinner Reservations API"))
	})
	http.HandleFunc("/reservations", getReservations)
	http.HandleFunc("/reservation", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createReservation(w, r)
		case http.MethodGet:
			getReservation(w, r)
		case http.MethodPut:
			updateReservation(w, r)
		case http.MethodDelete:
			deleteReservation(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
