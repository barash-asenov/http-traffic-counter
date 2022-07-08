package main

import (
	"reflect"
	"testing"
	"time"
)

func TestRequests_CountWithin(t *testing.T) {
	// create requests, older will be first
	requests := Requests{
		Request{time.Now().Add(-1 * time.Minute)},
		Request{time.Now().Add(-59 * time.Second)},
		Request{time.Now().Add(-5 * time.Second)},
		Request{time.Now()},
	}

	requestCount := requests.CountWithin(time.Now().Add(-1 * time.Minute))

	if requestCount != 3 {
		t.Errorf("expected 3, returned %d\n", requestCount)
	}

	requestCount = requests.CountWithin(time.Now().Add(-5 * time.Minute))

	if requestCount != 4 {
		t.Errorf("expected 4, returned %d\n", requestCount)
	}

	requestCount = requests.CountWithin(time.Now().Add(-1 * time.Second))

	if requestCount != 1 {
		t.Errorf("expected 1, returned %d\n", requestCount)
	}
}

func TestRequests_RemoveOlderFrom(t *testing.T) {
	// create requests, older will be first
	requests := Requests{
		Request{time.Now().Add(-1 * time.Minute)},
		Request{time.Now().Add(-59 * time.Second)},
		Request{time.Now().Add(-5 * time.Second)},
		Request{time.Now()},
	}

	requests.RemoveOlderFrom(time.Now())

	if !reflect.DeepEqual(requests, Requests{}) {
		t.Error("expected to remove all of them")
	}
}

func TestRequests_Add(t *testing.T) {
	requests := Requests{}
	newRequest := &Request{time.Now()}

	requests.Add(newRequest)

	if !reflect.DeepEqual(Requests{*newRequest}, requests) {
		t.Error("not successful add operation")
	}
}
