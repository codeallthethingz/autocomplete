package main

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/searchspring/autocomplete/process-data/sshttp"
)

var dataLocation string = ""
var communicationChannel chan string = make(chan string, 2)

func main() {
	routes := defineEndpoints()
	adminRoutes, err := sshttp.AdminEndpoints([]string{"prometheus", "metrics", "health", "version", "debug/pprof"})
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	port, adminPort := sshttp.Ports()

	if err := setupGlobalConfig(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	sshttp.StartServer("autocomplete data", routes, port, adminRoutes, adminPort)

}

func setupGlobalConfig() error {
	dataLocation = os.Getenv("DATA_LOCATION")
	if dataLocation == "" {
		return fmt.Errorf("DATA_LOCATION environment variable is required")
	}
	return nil
}

// return gorilla mux endpoints
func defineEndpoints() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/data/{siteId}", dataHandler).Methods("POST")
	r.HandleFunc("/autocomplete/{siteId}", autocompleteHandler).Methods("GET")
	return r
}
