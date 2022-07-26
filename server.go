package main

import (
	"fmt"
	"net/http"
	"time"
)

func RequestCountHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	newRequest := &Request{Timestamp: time.Now()}

	// add new request to object
	requests.Add(newRequest)

	time.Sleep(5 * time.Second)

	// write requests count withing moving window to response-writer
	requestCount := requests.CountWithin(currentTime.Add(-serverConfig.MovingWindow))
	requests.mu.Lock()
	defer requests.mu.Unlock()
	fmt.Fprintf(w, "Total requests within moving window: %d\n", requestCount)
}
