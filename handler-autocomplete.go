package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/searchspring/autocomplete/stringsearch"
)

var siteID2AutocompleteTrie = make(map[string]*stringsearch.AutocompleteTrie)

func autocompleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["siteId"]
	// read q query parameter
	q := r.URL.Query().Get("q")
	if q == "" {
		emptyResponse(w)
		return
	}
	var autocomplete *stringsearch.AutocompleteTrie
	if _, ok := siteID2AutocompleteTrie[siteID]; !ok {
		a, err := loadAutocomplete(siteID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		siteID2AutocompleteTrie[siteID] = a
	}
	if _, ok := siteID2AutocompleteTrie[siteID]; !ok {
		emptyResponse(w)
		return
	}

	autocomplete = siteID2AutocompleteTrie[siteID]

	if responses, ok := autocomplete.FindCaseAware(q); ok {
		// wrap response in object
		wrapped := map[string]interface{}{"data": responses}
		// convert to json

		json, err := json.Marshal(wrapped)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	} else {
		emptyResponse(w)
	}
}

func loadAutocomplete(siteID string) (*stringsearch.AutocompleteTrie, error) {
	// read the data file from disk
	dataFile := dataLocation + "/" + siteID
	f, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// create reader from file
	reader := bufio.NewReader(f)
	autocomplete := stringsearch.NewAutocompleteTrie(reader, 5)
	return autocomplete, nil
}

func emptyResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func reloadData(siteID string) {
	autocomplete, err := loadAutocomplete(siteID)
	if err != nil {
		return
	}
	siteID2AutocompleteTrie[siteID] = autocomplete
}
