package main

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/searchspring/autocomplete/process-data/disksort"
	"github.com/searchspring/autocomplete/process-data/ngramify"
)

func dataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["siteId"]

	// read the body contents to a string
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	// ngramify the body
	ngramifier, err := ngramify.New(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ngramifier.Ngramify(string(body), 3)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// sort and dedup the ngrams
	// get lines from the buffer
	lines := strings.Split(buf.String(), "\n")
	stringChan := make(chan string, 2)
	go func() {
		for _, line := range lines {
			stringChan <- line
		}
		close(stringChan)
	}()

	// create data file
	dataFile := dataLocation + "/" + siteID + ".txt"
	f, err := os.Create(dataFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	err = disksort.Sort(stringChan, writer, true)
	writer.Flush()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send success code
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`"{ "status": "success", "message": "data processed" }"`))
	communicationChannel <- siteID

}
