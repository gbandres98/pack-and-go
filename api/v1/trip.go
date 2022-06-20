package api_v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gbandres98/pack-and-go/db"
	"github.com/gbandres98/pack-and-go/model"
	"github.com/gorilla/mux"
)

type tripService interface {
	GetAllTrips() []model.Trip
	GetTripById(int32) (model.Trip, error)
	AddTrip(model.Trip) (model.Trip, error)
	GetTripPretty(model.Trip) (model.TripPretty, error)
}

type tripController struct {
	tripService
}

func NewTripController(tripService tripService) *tripController {
	return &tripController{tripService}
}

func (tripController *tripController) GetAllTrips(w http.ResponseWriter, req *http.Request) {
	trips := tripController.tripService.GetAllTrips()
	tripsPretty := []model.TripPretty{}

	for _, trip := range trips {
		tripPretty, err := tripController.tripService.GetTripPretty(trip)
		if err != nil {
			http.Error(w, fmt.Sprintf("Internal Server Error - %v", err), http.StatusInternalServerError)
			return
		}

		tripsPretty = append(tripsPretty, tripPretty)
	}

	body, _ := json.Marshal(tripsPretty)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (tripController *tripController) GetTripById(w http.ResponseWriter, req *http.Request) {
	idVar := mux.Vars(req)["id"]
	id, err := strconv.ParseInt(idVar, 10, 32)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad Request - invalid trip id: %v", err), http.StatusBadRequest)
		return
	}

	trip, err := tripController.tripService.GetTripById(int32(id))
	if err == db.ErrorTripNotFound {
		http.Error(w, fmt.Sprintf("Not Found - no trip found with id: %v", id), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error - %v", err), http.StatusInternalServerError)
		return
	}

	tripPretty, err := tripController.tripService.GetTripPretty(trip)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error - %v", err), http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(tripPretty)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (tripController *tripController) AddTrip(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)

	var newTrip model.Trip
	err := json.Unmarshal(requestBody, &newTrip)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad Request - invalid trip json: %v", err), http.StatusBadRequest)
		return
	}

	savedTrip, err := tripController.tripService.AddTrip(newTrip)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad Request - invalid trip: %v", err), http.StatusBadRequest)
		return
	}

	tripPretty, err := tripController.tripService.GetTripPretty(savedTrip)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error - %v", err), http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(tripPretty)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}