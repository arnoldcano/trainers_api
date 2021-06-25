package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Trainer struct {
	ID        int    `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func GetTrainerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	q := `
		SELECT id, email, phone, first_name, last_name
		FROM trainers
		WHERE id = $1
	`
	var t Trainer
	row := db.QueryRow(q, vars["id"])
	if err := row.Scan(
		&t.ID,
		&t.Email,
		&t.Phone,
		&t.FirstName,
		&t.LastName,
	); err != nil && err != sql.ErrNoRows {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func CreateTrainerHandler(w http.ResponseWriter, r *http.Request) {
	var t Trainer
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	q := `
		INSERT INTO trainers (email, phone, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	row := db.QueryRow(q, t.Email, t.Phone, t.FirstName, t.LastName)
	if err := row.Scan(&t.ID); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	var err error
	databaseURL := fmt.Sprintf(
		"postgres://%s@%s:5432/%s?sslmode=disable",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_NAME"),
	)
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/trainers/{id}", GetTrainerHandler).Methods("GET")
	router.HandleFunc("/trainers", CreateTrainerHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
