package api_v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gbandres98/pack-and-go/db"
	"github.com/gbandres98/pack-and-go/model"
	"github.com/gorilla/mux"
)

var testTrips = []model.Trip{
	{Id: 1, OriginId: 1, DestinationId: 2, Dates: "Mon Tue Wed Fri", Price: 40.55},
	{Id: 2, OriginId: 2, DestinationId: 1, Dates: "Sat Sun", Price: 40.55},
}

var testTripPretty = model.TripPretty{Id: 1, Origin: "Sevilla", Destination: "Madrid", Dates: "Mon Tue", Price: 40.55}

type mockTripService struct{
	failGetTripPretty bool
	failGetTripById bool
}

func (mockTripService *mockTripService) GetAllTrips() []model.Trip {
	return testTrips
}

func (mockTripService *mockTripService) GetTripById(id int32) (model.Trip, error) {
	if (mockTripService.failGetTripById) {
		return model.Trip{}, fmt.Errorf("test error")
	}

	if (id < 3) {
		return testTrips[id - 1], nil
	}

	return model.Trip{}, db.ErrorTripNotFound
}

func (mockTripService *mockTripService) AddTrip(trip model.Trip) (model.Trip, error) {
	if trip.OriginId > 2 || trip.DestinationId > 2 || trip.OriginId < 1 || trip.DestinationId < 1 {
		return model.Trip{}, errors.New("invalid originId or destinationId")
	}
	trip.Id = 3
	return trip, nil
}

func (mockTripService *mockTripService) GetTripPretty(trip model.Trip) (model.TripPretty, error) {
	if mockTripService.failGetTripPretty {
		return model.TripPretty{}, fmt.Errorf("test error")
	}
	result := testTripPretty
	result.Id = trip.Id
	return result, nil
}

func TestGetAllTrips_1(t *testing.T) {
	tripController := NewTripController(&mockTripService{})

	req := httptest.NewRequest("GET", "/trip", nil)
	responseRecorder := httptest.NewRecorder()

	tripController.GetAllTrips(responseRecorder, req)

	expected := strings.TrimSpace(
		`[{"id":1,"origin":"Sevilla","destination":"Madrid","dates":"Mon Tue","price":40.55},` + 
		`{"id":2,"origin":"Sevilla","destination":"Madrid","dates":"Mon Tue","price":40.55}]`)
	result := strings.TrimSpace(responseRecorder.Body.String())

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected response code to be %v, got %v", http.StatusOK, responseRecorder.Code)
	}
	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestGetAllTrips_2(t *testing.T) {
	tripController := NewTripController(&mockTripService{ failGetTripPretty: true })

	req := httptest.NewRequest("GET", "/trip", nil)
	responseRecorder := httptest.NewRecorder()

	tripController.GetAllTrips(responseRecorder, req)

	if responseRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected response code to be %v, got %v", http.StatusInternalServerError, responseRecorder.Code)
	}
}

func TestGetTripById_1(t *testing.T) {
	tripController := NewTripController(&mockTripService{})

	req := httptest.NewRequest("GET", "/trip/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	responseRecorder := httptest.NewRecorder()

	tripController.GetTripById(responseRecorder, req)

	expected := strings.TrimSpace(
		`{"id":1,"origin":"Sevilla","destination":"Madrid","dates":"Mon Tue","price":40.55}`)
	result := strings.TrimSpace(responseRecorder.Body.String())

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected response code to be %v, got %v", http.StatusOK, responseRecorder.Code)
	}
	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestGetTripById_2(t *testing.T) {
	tripController := NewTripController(&mockTripService{})

	req := httptest.NewRequest("GET", "/trip/3", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "3"})
	responseRecorder := httptest.NewRecorder()

	tripController.GetTripById(responseRecorder, req)

	if responseRecorder.Code != http.StatusNotFound {
		t.Fatalf("expected response code to be %v, got %v", http.StatusNotFound, responseRecorder.Code)
	}
}

func TestGetTripById_3(t *testing.T) {
	tripController := NewTripController(&mockTripService{ failGetTripById: true })

	req := httptest.NewRequest("GET", "/trip/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	responseRecorder := httptest.NewRecorder()

	tripController.GetTripById(responseRecorder, req)

	if responseRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected response code to be %v, got %v", http.StatusInternalServerError, responseRecorder.Code)
	}
}

func TestGetTripById_4(t *testing.T) {
	tripController := NewTripController(&mockTripService{ failGetTripPretty: true })

	req := httptest.NewRequest("GET", "/trip/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	responseRecorder := httptest.NewRecorder()

	tripController.GetTripById(responseRecorder, req)

	if responseRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected response code to be %v, got %v", http.StatusInternalServerError, responseRecorder.Code)
	}
}

func TestGetTripById_5(t *testing.T) {
	tripController := NewTripController(&mockTripService{})

	req := httptest.NewRequest("GET", "/trip/a", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "a"})
	responseRecorder := httptest.NewRecorder()

	tripController.GetTripById(responseRecorder, req)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected response code to be %v, got %v", http.StatusBadRequest, responseRecorder.Code)
	}
}

func TestAddTrip_1(t *testing.T) {
	tripController := NewTripController(&mockTripService{})

	req := httptest.NewRequest("POST", "/trip", strings.NewReader(`{"originId":1,"destinationId":2,"dates":"Mon Tue","price":40.55}`))
	responseRecorder := httptest.NewRecorder()

	tripController.AddTrip(responseRecorder, req)

	expected := strings.TrimSpace(
		`{"id":3,"origin":"Sevilla","destination":"Madrid","dates":"Mon Tue","price":40.55}`)
	result := strings.TrimSpace(responseRecorder.Body.String())

	if responseRecorder.Code != http.StatusCreated {
		t.Fatalf("expected response code to be %v, got %v", http.StatusCreated, responseRecorder.Code)
	}
	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestAddTrip_2(t *testing.T) {
	tripController := NewTripController(&mockTripService{})

	req := httptest.NewRequest("POST", "/trip", strings.NewReader(`{"id":5,"originId":1,"destinationId":2,"dates":"Mon Tue","price":40.55}`))
	responseRecorder := httptest.NewRecorder()

	tripController.AddTrip(responseRecorder, req)

	expected := strings.TrimSpace(
		`{"id":3,"origin":"Sevilla","destination":"Madrid","dates":"Mon Tue","price":40.55}`)
	result := strings.TrimSpace(responseRecorder.Body.String())

	if responseRecorder.Code != http.StatusCreated {
		t.Fatalf("expected response code to be %v, got %v", http.StatusCreated, responseRecorder.Code)
	}
	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestAddTrip_3(t *testing.T) {
	tripController := NewTripController(&mockTripService{})

	req := httptest.NewRequest("POST", "/trip", nil)
	responseRecorder := httptest.NewRecorder()

	tripController.AddTrip(responseRecorder, req)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected response code to be %v, got %v", http.StatusBadRequest, responseRecorder.Code)
	}
}

func TestAddTrip_4(t *testing.T) {
	tripController := NewTripController(&mockTripService{ failGetTripPretty: true})

	req := httptest.NewRequest("POST", "/trip", strings.NewReader(`{"id":5,"originId":1,"destinationId":2,"dates":"Mon Tue","price":40.55}`))
	responseRecorder := httptest.NewRecorder()

	tripController.AddTrip(responseRecorder, req)

	if responseRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected response code to be %v, got %v", http.StatusInternalServerError, responseRecorder.Code)
	}
}

func TestAddTrip_5(t *testing.T) {
	tripController := NewTripController(&mockTripService{})

	req := httptest.NewRequest("POST", "/trip", strings.NewReader(`{"id":5,"originId":1,"destinationId":3,"dates":"Mon Tue","price":40.55}`))
	responseRecorder := httptest.NewRecorder()

	tripController.AddTrip(responseRecorder, req)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected response code to be %v, got %v", http.StatusBadRequest, responseRecorder.Code)
	}
}