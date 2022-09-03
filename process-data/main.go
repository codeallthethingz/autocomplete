package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/searchspring/autocomplete/process-data/sshttp"
)

func main() {
	routes := defineEndpoints()
	adminRoutes, err := sshttp.AdminEndpoints([]string{"prometheus", "metrics", "health", "version", "debug/pprof"})
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	port, adminPort := ports()
	sshttp.StartServer("autocomplete Data", routes, port, adminRoutes, adminPort)

}

func ports() (string, string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	adminPort := os.Getenv("ADMIN_PORT")
	if adminPort == "" {
		adminPort = "8081"
	}
	return port, adminPort
}

// return gorilla mux endpoints
func defineEndpoints() *mux.Router {
	r := mux.NewRouter()
	defineEndpointData(r)
	return r
}

// defineEndpointData defines the data endpoint
func defineEndpointData(r *mux.Router) {
	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
}
