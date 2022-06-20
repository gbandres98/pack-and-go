package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gbandres98/pack-and-go/model"
)

func TestGetAllTrips(t *testing.T) {
	app := setupApplication(applicationConfig{
		fileDBPath: "./cities_test.txt",
	})

	req := httptest.NewRequest("GET", "/api/v1/trip", nil)
	responseRecorder := httptest.NewRecorder()

	app.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected response code to be %v, got %v", http.StatusOK, responseRecorder.Code)
	}
}

func TestAddAndGetTrip(t *testing.T) {
	app := setupApplication(applicationConfig{
		fileDBPath: "./cities_test.txt",
	})

	req := httptest.NewRequest("POST", "/api/v1/trip", strings.NewReader(`{"originId":1,"destinationId":2,"dates":"Mon Tue","price":40.55}`))
	responseRecorder := httptest.NewRecorder()

	app.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusCreated {
		t.Fatalf("expected response code to be %v, got %v", http.StatusOK, responseRecorder.Code)
	}
	
	var createdTrip model.TripPretty
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &createdTrip)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/trip/%v", createdTrip.Id), nil)
	responseRecorder = httptest.NewRecorder()

	app.ServeHTTP(responseRecorder, req)

	expected := strings.TrimSpace(
		`{"id":` + strconv.Itoa(int(createdTrip.Id)) + `,"origin":"Barcelona","destination":"Seville","dates":"Mon Tue","price":40.55}`)
	result := strings.TrimSpace(responseRecorder.Body.String())

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected response code to be %v, got %v", http.StatusOK, responseRecorder.Code)
	}
	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestAddAndGetAllTrips(t *testing.T) {
	app := setupApplication(applicationConfig{
		fileDBPath: "./cities_test.txt",
	})

	req := httptest.NewRequest("POST", "/api/v1/trip", strings.NewReader(`{"originId":1,"destinationId":2,"dates":"Mon Tue","price":40.55}`))
	responseRecorder := httptest.NewRecorder()

	app.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusCreated {
		t.Fatalf("expected response code to be %v, got %v", http.StatusOK, responseRecorder.Code)
	}
	
	var createdTrip model.TripPretty
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &createdTrip)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req = httptest.NewRequest("GET", "/api/v1/trip", nil)
	responseRecorder = httptest.NewRecorder()

	app.ServeHTTP(responseRecorder, req)

	expected := strings.TrimSpace(
		`{"id":` + strconv.Itoa(int(createdTrip.Id)) + `,"origin":"Barcelona","destination":"Seville","dates":"Mon Tue","price":40.55}`)
	result := strings.TrimSpace(responseRecorder.Body.String())

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected response code to be %v, got %v", http.StatusOK, responseRecorder.Code)
	}
	if !strings.Contains(result, expected) {
		t.Fatalf("expected %v to include %v", result, expected)
	}
}