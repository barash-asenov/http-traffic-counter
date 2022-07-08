package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRequestCountHandler(t *testing.T) {
	serverConfig = &ServerConfig{
		ExportFileName: "test-requests.json",
		MovingWindow:   5 * time.Second,
	}

	t.Run("returns correct response on single executions", func(t *testing.T) {
		requests = &Requests{}

		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Error(err)
		}
		response := httptest.NewRecorder()

		RequestCountHandler(response, request)

		res := response.Body.String()
		expected := "Total requests within moving window: 1\n"

		if res != expected {
			t.Errorf("expected 1, returned %s\n", res)
		}
	})

	t.Run("returns correct response on multiple executions", func(t *testing.T) {
		requests = &Requests{}

		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Error(err)
		}
		response := httptest.NewRecorder()

		var wg sync.WaitGroup

		wg.Add(1000)
		for i := 0; i < 1000; i++ {
			func() {
				defer wg.Done()

				RequestCountHandler(response, request)
			}()
		}

		wg.Wait()

		response = httptest.NewRecorder()
		RequestCountHandler(response, request)

		res := response.Body.String()
		expected := "Total requests within moving window: 1001\n"

		if res != expected {
			t.Errorf("expected 1001, returned %s\n", res)
		}
	})

	t.Cleanup(clearTestResources)
}

// remove the created file
func clearTestResources() {
	if _, err := os.Stat(serverConfig.ExportFileName); err == nil {
		// remove export file...
		// normally not very nice to create and delete file that is used in main program... but for
		// the sake of simplicity
		err := os.Remove(serverConfig.ExportFileName)
		if err != nil {
			log.Println("Failed to remove test file!")
		}
	}
}
