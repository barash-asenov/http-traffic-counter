package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const DefaultExportFileName = "requests.json"
const DefaultMovingWindow = 1 * time.Minute

type ServerConfig struct {
	ExportFileName string
	MovingWindow   time.Duration
}

var requests = &Requests{Data: []Request{}}
var serverConfig = &ServerConfig{
	ExportFileName: DefaultExportFileName,
	MovingWindow:   DefaultMovingWindow,
}

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
