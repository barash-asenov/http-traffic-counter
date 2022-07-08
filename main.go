package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const ExportFileName = "data/requests.json"
const MovingWindow = 1 * time.Minute

func main() {
	// load file from disk to ram
	requests, err := initializeRequests()

	if err != nil {
		log.Fatalf("failed to load requests from disk: %s\n", err)
	}

	// main handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		newRequest := Request{Timestamp: time.Now()}

		// add new request to object
		requests.Add(newRequest)

		// clean requests older than moving window and persist the requests in goroutine
		go func() {
			requests.RemoveOlderFrom(time.Now().Add(-MovingWindow))

			// persist the requests in disk...
			err := exportRequests(requests)
			if err != nil {
				log.Printf("export failed %s\n", err)
			}
		}()

		// write requests count withing moving window to response-writer
		requestCount := requests.CountWithin(time.Now().Add(-MovingWindow))

		fmt.Fprintf(w, "Total requests: %d\n", requestCount)
	})

	// start web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func initializeRequests() (*Requests, error) {
	content, err := ioutil.ReadFile(ExportFileName)

	if errors.Is(err, os.ErrNotExist) {
		// file not exists yet... initialize it empty
		return &Requests{}, nil
	} else if err != nil {
		// not known error...
		return nil, err
	} else {
		return LoadRequests(content)
	}
}

func exportRequests(r *Requests) error {
	data, err := r.AsJSON()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(ExportFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
