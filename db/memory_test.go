package db

import (
	"reflect"
	"testing"

	"github.com/gbandres98/pack-and-go/model"
)

var testTrips = []model.Trip{
	{Id: 1, OriginId: 1, DestinationId: 2, Dates: "Mon Tue Wed Fri", Price: 40.55},
	{Id: 2, OriginId: 2, DestinationId: 1, Dates: "Sat Sun", Price: 40.55},
}

var newTrip = model.Trip{OriginId: 3, DestinationId: 6, Dates: "Mon Tue Wed Thu Fri", Price: 32.10}
var newTripWithId = model.Trip{OriginId: 1, DestinationId: 6, Dates: "Mon Tue Wed Thu Fri", Price: 32.10}

func TestNewMemoryDB(t *testing.T) {
	memoryDB := NewMemoryDB()

	if (memoryDB.trips == nil) {
		t.Fatalf("expected non-nil array of trips")
	}
}

func TestGetAllTrips_1(t *testing.T) {
	memoryDB := memoryDB{ trips: testTrips }

	trips := memoryDB.GetAllTrips()
	if !reflect.DeepEqual(trips, testTrips) {
		t.Fatalf("expected %v, got %v", testTrips, trips)
	}
}

func TestGetAllTrips_2(t *testing.T) {
	defer func() { recover() }()

	memoryDB := memoryDB{}
	memoryDB.GetAllTrips()

	t.Fatalf("expected GetAllTrips() to panic")
}

func TestGetTripById_1(t *testing.T) {
	memoryDB := memoryDB{ trips: testTrips }

	trip, err := memoryDB.GetTripById(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(trip, testTrips[0]) {
		t.Fatalf("expected %v, got %v", testTrips, trip)
	}
}

func TestGetTripById_2(t *testing.T) {
	memoryDB := memoryDB{ trips: testTrips }

	_, err := memoryDB.GetTripById(3)
	if err != ErrorTripNotFound {
		t.Fatalf("expected error: %v, got error: %v", ErrorTripNotFound, err)
	}
}

func TestGetTripById_3(t *testing.T) {
	defer func() { recover() }()

	memoryDB := memoryDB{}
	memoryDB.GetTripById(3)

	t.Fatalf("expected GetTripById() to panic")
}

func TestAddTrip_1(t *testing.T) {
	memoryDB := memoryDB{ trips: testTrips, nextId: 3 }

	savedTrip := memoryDB.AddTrip(newTrip)
	if savedTrip.Id != 3 {
		t.Fatalf("expected new trip to have id %v, got id %v", 3, savedTrip.Id)
	}

	trips := memoryDB.GetAllTrips()
	if len(trips) != 3 {
		t.Fatalf("expected trip list to have length %v, got %v", 3, len(trips))
	}
}

func TestAddTrip_2(t *testing.T) {
	memoryDB := memoryDB{ trips: testTrips, nextId: 3 }

	savedTrip := memoryDB.AddTrip(newTripWithId)
	if savedTrip.Id != 3 {
		t.Fatalf("expected new trip to have id %v, got id %v", 3, savedTrip.Id)
	}

	trips := memoryDB.GetAllTrips()
	if len(trips) != 3 {
		t.Fatalf("expected trip list to have length %v, got %v", 3, len(trips))
	}
}