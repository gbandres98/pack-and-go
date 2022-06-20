package service

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gbandres98/pack-and-go/model"
)

type cityDB interface {
	GetCityById(int32) (model.City, error)
}

type tripDB interface {
	GetAllTrips() []model.Trip
	GetTripById(int32) (model.Trip, error)
	AddTrip(model.Trip) model.Trip
}

type tripService struct {
	cityDB
	tripDB
}

var datesRegexp *regexp.Regexp

func init() {
	var err error
	datesRegexp, err = regexp.Compile(`^([A-Z][a-z][a-z] )*([A-Z][a-z][a-z])$`)
	if err != nil {
		log.Panicln(err)
	}
}

func NewTripService(cityDB cityDB, tripDB tripDB) *tripService {
	return &tripService{cityDB, tripDB}
}

func (tripService *tripService) GetAllTrips() []model.Trip {
	return tripService.tripDB.GetAllTrips()
}

func (tripService *tripService) GetTripById(id int32) (model.Trip, error) {
	return tripService.tripDB.GetTripById(id)
}

func (tripService *tripService) AddTrip(trip model.Trip) (model.Trip, error) {
	if !datesRegexp.MatchString(trip.Dates) {
		return model.Trip{}, fmt.Errorf("invalid dates format: %v", trip.Dates)
	}

	_, err := tripService.cityDB.GetCityById(trip.OriginId)
	if (err != nil) {
		return model.Trip{}, fmt.Errorf("could not find origin city with id: %v", trip.OriginId)
	}

	_, err = tripService.cityDB.GetCityById(trip.DestinationId)
	if (err != nil) {
		return model.Trip{}, fmt.Errorf("could not find destination city with id: %v", trip.DestinationId)
	}

	return tripService.tripDB.AddTrip(trip), nil
}

func (tripService *tripService) GetTripPretty(trip model.Trip) (model.TripPretty, error) {
	originCity, err := tripService.cityDB.GetCityById(trip.OriginId)
	if (err != nil) {
		return model.TripPretty{}, fmt.Errorf("could not find origin city with id: %v", trip.OriginId)
	}

	destinationCity, err := tripService.cityDB.GetCityById(trip.DestinationId)
	if (err != nil) {
		return model.TripPretty{}, fmt.Errorf("could not find destination city with id: %v", trip.DestinationId)
	}

	tripPretty := model.TripPretty{
		Id: trip.Id,
		Origin: originCity.Name,
		Destination: destinationCity.Name,
		Dates: trip.Dates,
		Price: trip.Price,
	}

	return tripPretty, nil
}