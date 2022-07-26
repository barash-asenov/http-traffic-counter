package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const DefaultExportFileName = "requests.json"
const DefaultMovingWindow = 1 * time.Minute
const MaxRequests = 5

type ServerConfig struct {
	ExportFileName string
	MovingWindow   time.Duration
	RequestLimit   chan struct{}
}

var requests = &Requests{Data: []Request{}}
var serverConfig = &ServerConfig{
	ExportFileName: DefaultExportFileName,
	MovingWindow:   DefaultMovingWindow,
	RequestLimit:   make(chan struct{}, MaxRequests),
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// load file from disk to ram
	err := initializeRequests()
	if err != nil {
		log.Fatalf("failed to load requests from disk: %s\n", err)
	}

	// main handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", RequestCountHandler)

	jobsFinished := make(chan struct{})

	go func() {
		defer close(jobsFinished)

		for {
			select {
			case <-ctx.Done():
				// os interrupt or sigterm
				exportRequests()
				jobsFinished <- struct{}{}
				return
			case <-time.After(10 * time.Second):
				clearRequests()
				exportRequests()
			}
		}
	}()

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: limit(mux),
	}

	go func() {
		err := httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("shutting down server...")
		} else if err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	<-jobsFinished
	err = httpServer.Shutdown(context.Background())

	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("good bye...")
	}
}

func clearRequests() {
	currentTime := time.Now()
	removedRequests := requests.RemoveOlderFrom(currentTime.Add(-serverConfig.MovingWindow))
	log.Printf("cleared %d requests\n", removedRequests)
}

func exportRequests() {
	data, err := requests.AsJSON()
	if err != nil {
		log.Printf("export failed %s\n", err)
	}

	err = ioutil.WriteFile(serverConfig.ExportFileName, data, 0644)
	if err != nil {
		log.Printf("export failed %s\n", err)
	}

	log.Println("exported successfully!")
}

func initializeRequests() error {
	content, err := ioutil.ReadFile(serverConfig.ExportFileName)

	if errors.Is(err, os.ErrNotExist) {
		// file not exists yet... initialize it empty
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

	requests.Data = unmarshaledRequests

	return nil
}
