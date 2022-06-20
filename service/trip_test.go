package service

import (
	"reflect"
	"testing"

	"github.com/gbandres98/pack-and-go/db"
	"github.com/gbandres98/pack-and-go/model"
)

var testCities = []model.City{
	{Id: 1, Name: "Sevilla"},
	{Id: 2, Name: "Madrid"},
}

var testTrips = []model.Trip{
	{Id: 1, OriginId: 1, DestinationId: 2, Dates: "Mon Tue Wed Fri", Price: 40.55},
	{Id: 2, OriginId: 2, DestinationId: 1, Dates: "Sat Sun", Price: 40.55},
}

type mockCityDB struct{}

type mockTripDB struct{}

func (mockCityDB *mockCityDB) GetCityById(id int32) (model.City, error) {
	if (id < 3) {
		return testCities[id - 1], nil
	}

	return model.City{}, db.ErrorCityNotFound
}

func (mockTripDB *mockTripDB) GetAllTrips() []model.Trip {
	return testTrips
}

func (mockTripDB *mockTripDB) GetTripById(id int32) (model.Trip, error) {
	if (id < 3) {
		return testTrips[id - 1], nil
	}

	return model.Trip{}, db.ErrorTripNotFound
}

func (mockTripDB *mockTripDB) AddTrip(trip model.Trip) model.Trip {
	trip.Id = 3
	return trip
}

func TestGetAllTrips_1(t *testing.T) {
	tripService := NewTripService(&mockCityDB{}, &mockTripDB{})

	trips := tripService.GetAllTrips()
	if !reflect.DeepEqual(trips, testTrips) {
		t.Fatalf("expected %v, got %v", testTrips, trips)
	}
}

func TestGetTripById_1(t *testing.T) {
	tripService := NewTripService(&mockCityDB{}, &mockTripDB{})

	trip, err := tripService.GetTripById(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(trip, testTrips[0]) {
		t.Fatalf("expected %v, got %v", testTrips, trip)
	}
}

func TestGetTripById_2(t *testing.T) {
	tripService := NewTripService(&mockCityDB{}, &mockTripDB{})

	_, err := tripService.GetTripById(3)
	if err != db.ErrorTripNotFound {
		t.Fatalf("expected error: %v, got error: %v", db.ErrorTripNotFound, err)
	}
}

func TestAddTrip_1(t *testing.T) {
	tripService := NewTripService(&mockCityDB{}, &mockTripDB{})

	newTrip := model.Trip{OriginId: 1, DestinationId: 2, Dates: "Mon Tue", Price: 40.21}

	savedTrip, err := tripService.AddTrip(newTrip)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if savedTrip.Id != 3 {
		t.Fatalf("expected saved trip to have id: %v, got %v", 3, savedTrip)
	}
}

func TestAddTrip_2(t *testing.T) {
	tripService := NewTripService(&mockCityDB{}, &mockTripDB{})

	newTrip := model.Trip{OriginId: 3, DestinationId: 2, Dates: "Mon Tue", Price: 40.21}

	_, err := tripService.AddTrip(newTrip)
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddTrip_3(t *testing.T) {
	tripService := NewTripService(&mockCityDB{}, &mockTripDB{})

	newTrip := model.Trip{OriginId: 1, DestinationId: 3, Dates: "Mon Tue", Price: 40.21}

	_, err := tripService.AddTrip(newTrip)
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddTrip_4(t *testing.T) {
	tripService := NewTripService(&mockCityDB{}, &mockTripDB{})

	newTrip := model.Trip{OriginId: 1, DestinationId: 2, Dates: "MonTue", Price: 40.21}

	_, err := tripService.AddTrip(newTrip)
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddTrip_5(t *testing.T) {
	tripService := NewTripService(&mockCityDB{}, &mockTripDB{})

	newTrip := model.Trip{OriginId: 1, DestinationId: 2, Dates: "mon Tue", Price: 40.21}

	_, err := tripService.AddTrip(newTrip)
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestGetTripPretty_1(t *testing.T) {
	tripService := NewTripService(&mockCityDB{}, &mockTripDB{})

	trip := model.Trip{Id: 3, OriginId: 1, DestinationId: 2, Dates: "Mon Tue", Price: 40.21}
	expected := model.TripPretty{Id: 3, Origin: "Sevilla", Destination: "Madrid", Dates: "Mon Tue", Price: 40.21}

	tripPretty, err := tripService.GetTripPretty(trip)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(tripPretty, expected) {
		t.Fatalf("expected %v, got %v", expected, tripPretty)
	}
}

func TestGetTripPretty_2(t *testing.T) {
	tripService := NewTripService(&mockCityDB{}, &mockTripDB{})

	trip := model.Trip{Id: 3, OriginId: 3, DestinationId: 2, Dates: "Mon Tue", Price: 40.21}

	_, err := tripService.GetTripPretty(trip)
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestGetTripPretty_3(t *testing.T) {
	tripService := NewTripService(&mockCityDB{}, &mockTripDB{})

	trip := model.Trip{Id: 3, OriginId: 2, DestinationId: 3, Dates: "Mon Tue", Price: 40.21}

	_, err := tripService.GetTripPretty(trip)
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
}