package api_v1

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetRoutes(router *mux.Router, tripController tripController) *mux.Router {
	router.StrictSlash(true)

	router.HandleFunc("/trip", tripController.GetAllTrips).Methods(http.MethodGet)
	router.HandleFunc("/trip", tripController.AddTrip).Methods(http.MethodPost)
	router.HandleFunc("/trip/{id}", tripController.GetTripById).Methods(http.MethodGet)

	return router
}