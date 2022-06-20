package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	ip := *flag.String("ip", "", "IP Address for the application server to listen at")
	port := *flag.String("port", "8080", "Port for the application server to listen at")
	fileDBPath := *flag.String("db_file", "cities.txt", "Path to the file to be used as file DB")

	address := fmt.Sprintf("%v:%v", ip, port)

	server := http.Server{
		Addr: address,
		Handler: setupApplication(applicationConfig{
			fileDBPath: fileDBPath,
		}),
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("PackAndGo Server listening on %v", address)

	log.Panic(server.ListenAndServe())
}