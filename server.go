package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func RequestCountHandler(w http.ResponseWriter, r *http.Request) {
	newRequest := &Request{Timestamp: time.Now()}

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

	fmt.Fprintf(w, "Total requests within moving window: %d\n", requestCount)
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
