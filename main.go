package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/searchspring/autocomplete/sshttp"
)

var dataLocation string = ""

//go:embed index.html
var indexPage []byte

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
	// create data directory if it doesn't exist
	if _, err := os.Stat(dataLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(dataLocation, 0755); err != nil {
			return fmt.Errorf("error creating data directory: %v", err)
		}
	}

	return nil
}

// return gorilla mux endpoints
func defineEndpoints() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(indexPage)
	}).Methods("GET")
	r.HandleFunc("/data/{siteId}", dataHandler).Methods("POST")
	r.HandleFunc("/autocomplete/{siteId}", autocompleteHandler).Methods("GET")
	return r
}
