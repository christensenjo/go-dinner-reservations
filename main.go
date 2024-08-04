package main

import (
	"encoding/json"
	"log"
	"math/rand"
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
	mutex        sync.Mutex
)

func generateSampleReservations(num int) {
	names := []string{"John Doe", "Jane Smith", "Alice Johnson", "Bob Brown", "Emma Wilson"}
	times := []string{"18:00", "19:00", "20:00", "21:00"}
	dates := []string{"2024-09-01", "2024-09-02", "2024-09-03", "2024-09-04"}

	for i := 0; i < num; i++ {
		newReservation := Reservation{
			ID:     nextID,
			Name:   names[rand.Intn(len(names))],
			Date:   dates[rand.Intn(len(dates))],
			Time:   times[rand.Intn(len(times))],
			Guests: rand.Intn(10) + 1,
			Phone:  "555-0" + strconv.Itoa(rand.Intn(1000)+1000),
		}
		mutex.Lock()
		reservations = append(reservations, newReservation)
		nextID++
		mutex.Unlock()
	}
}

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
	// Populate the reservations slice with sample data
	generateSampleReservations(10)

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
