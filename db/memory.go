package db

import (
	"log"
	"sync"

	"github.com/gbandres98/pack-and-go/model"
)

var trips = []model.Trip{
	{Id: 1, OriginId: 1, DestinationId: 2, Dates: "Mon Tue Wed Fri", Price: 40.55},
	{Id: 2, OriginId: 2, DestinationId: 1, Dates: "Sat Sun", Price: 40.55},
	{Id: 3, OriginId: 3, DestinationId: 6, Dates: "Mon Tue Wed Thu Fri", Price: 32.10},
}

type memoryDB struct {
	trips []model.Trip
	nextId int32
	writeLock sync.Mutex
}

func NewMemoryDB() *memoryDB {
	return &memoryDB{trips: trips, nextId: 4}
}

func (memoryDB *memoryDB) GetAllTrips() ([]model.Trip) {
	if (memoryDB.trips == nil) {
		log.Panicln("non-initialized memory database")
	}

	return memoryDB.trips
}

func (memoryDB *memoryDB) GetTripById(id int32) (model.Trip, error) {
	if (memoryDB.trips == nil) {
		log.Panicln("non-initialized memory database")
	}

	for _, trip := range memoryDB.trips {
		if trip.Id == id {
			return trip, nil
		}
	}

	return model.Trip{}, ErrorTripNotFound
}

func (memoryDB *memoryDB) AddTrip(trip model.Trip) model.Trip {
	memoryDB.writeLock.Lock()
	defer memoryDB.writeLock.Unlock()

	trip.Id = memoryDB.nextId
	memoryDB.nextId++

	memoryDB.trips = append(memoryDB.trips, trip)
	return trip
}