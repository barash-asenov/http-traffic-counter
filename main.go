package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const ExportFileName = "requests.json"
const MovingWindow = 1 * time.Minute

var requests *Requests

func main() {
	// load file from disk to ram
	err := initializeRequests()
	if err != nil {
		log.Fatalf("failed to load requests from disk: %s\n", err)
	}

	// main handler
	http.HandleFunc("/", RequestCountHandler)

	// start web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func initializeRequests() error {
	content, err := ioutil.ReadFile(ExportFileName)

	if errors.Is(err, os.ErrNotExist) {
		// file not exists yet... initialize it empty
		requests = &Requests{}

		return nil
	} else if err != nil {
		// not known error...
		return err
	}

	// file exists.. unmarshal
	unmarshaledRequests, err := LoadRequests(content)
	if err != nil {
		return err
	}

	requests = unmarshaledRequests

	return nil
}
