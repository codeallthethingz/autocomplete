package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func autocompleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["siteId"]
	fmt.Println("siteId: ", siteID)
	// read q query parameter
	q := r.URL.Query().Get("q")
	if q == "" {
		emptyResponse(w)
	}
}

func emptyResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
