package main

import (
	"log"

	api_v1 "github.com/gbandres98/pack-and-go/api/v1"
	"github.com/gbandres98/pack-and-go/db"
	"github.com/gbandres98/pack-and-go/service"
	"github.com/gorilla/mux"
)

type applicationConfig struct {
	fileDBPath string
}

func setupApplication(applicationConfig applicationConfig) *mux.Router {
	// Databases
	fileDB := db.NewFileDB(applicationConfig.fileDBPath)
	memoryDB := db.NewMemoryDB()

	// Services
	tripService := service.NewTripService(fileDB, memoryDB)

	// Controllers
	tripController := api_v1.NewTripController(tripService)

	// Routes
	router := mux.NewRouter()	
	api_v1.SetRoutes(router.PathPrefix("/api/v1").Subrouter(), *tripController)

	logRoutes(router)
	
	return router
}

func logRoutes(router *mux.Router) {
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        path, _ := route.GetPathTemplate()
		method, _ := route.GetMethods()

		if (len(method) > 0) {
			log.Printf("%v %v", path, method)
		}        
        return nil
    })

}