package model

type Trip struct {
	Id            int32   `json:"id"`
	OriginId      int32   `json:"originId"`
	DestinationId int32   `json:"destinationId"`
	Dates         string  `json:"dates"`
	Price         float64 `json:"price"`
}

type TripPretty struct {
	Id          int32   `json:"id"`
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Dates       string  `json:"dates"`
	Price       float64 `json:"price"`
}

type City struct {
	Id   int32
	Name string
}