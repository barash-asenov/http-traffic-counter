package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRequestCountHandler(t *testing.T) {
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
}
