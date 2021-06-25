package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	databaseURL := fmt.Sprintf(
		"postgres://%s@%s:5432/%s?sslmode=disable",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_NAME"),
	)
	txdb.Register("txdb", "postgres", databaseURL)
	db, _ = sql.Open("txdb", databaseURL)
	p := filepath.Join("testdata", "seeds.sql")
	f, _ := ioutil.ReadFile(p)
	seeds := strings.Split(string(f), ";")
	for _, s := range seeds {
		db.Exec(s)
	}
	defer db.Close()

	os.Exit(m.Run())
}

func TestGetTrainer(t *testing.T) {
	r, err := http.NewRequest("GET", "/trainers/1", nil)
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/trainers/{id}", GetTrainerHandler).Methods("GET")
	router.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected %v, Got %v", http.StatusOK, status)
	}
	expected := `{"id":1,"email":"bruce@lee.com","phone":"11111111111","first_name":"bruce","last_name":"lee"}`
	if body := rr.Body.String(); reflect.DeepEqual(body, expected) {
		t.Errorf("Expected %v, Got %v", expected, body)
	}
}

func TestGetTrainerWithNoResults(t *testing.T) {
	r, err := http.NewRequest("GET", "/trainers/0", nil)
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/trainers/{id}", GetTrainerHandler).Methods("GET")
	router.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected %v, Got %v", http.StatusOK, status)
	}
	expected := "{}"
	if body := rr.Body.String(); reflect.DeepEqual(body, expected) {
		t.Errorf("Expected %v, Got %v", expected, body)
	}
}

func TestCreateTrainer(t *testing.T) {
	data := []byte(`{"email":"jackie@chan.com","phone":"33333333333","first_name":"jackie","last_name":"chan"}`)
	r, err := http.NewRequest("POST", "/trainers", bytes.NewBuffer(data))
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/trainers", CreateTrainerHandler).Methods("POST")
	router.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected %v, Got %v", http.StatusOK, status)
	}
}

func TestCreateTrainerWithBadRequest(t *testing.T) {
	data := []byte("boom")
	r, err := http.NewRequest("POST", "/trainers", bytes.NewBuffer(data))
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/trainers", CreateTrainerHandler).Methods("POST")
	router.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected %v, Got %v", http.StatusBadRequest, status)
	}
}
