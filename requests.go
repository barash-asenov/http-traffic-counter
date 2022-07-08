package main

import (
	"encoding/json"
	"sync"
	"time"
)

type Request struct {
	Timestamp time.Time `json:"timestamp"`
}

type Requests struct {
	mu   sync.Mutex
	Data []Request
}

// LoadRequests loads requests from disk
func LoadRequests(content []byte) ([]Request, error) {
	var requests []Request

	err := json.Unmarshal(content, &requests)

	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Requests) Add(request *Request) {
	r.mu.Lock()
	defer r.mu.Unlock()
	(*r).Data = append((*r).Data, *request)
}

func (r *Requests) CountWithin(time time.Time) int {
	totalCount := 0

	// iterate backwards, from the latest and the newest request
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := len((*r).Data) - 1; i >= 0; i-- {
		if (*r).Data[i].Timestamp.After(time) {
			totalCount++
		} else {
			// if it's less, no need to iterate more, all next values will be before the given time
			break
		}
	}

	return totalCount
}

// AsJSON returns byte array of marshalled JSON data
func (r *Requests) AsJSON() ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := json.Marshal((*r).Data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// RemoveOlderFrom remove records from the requests array, starting from the
// given timestamp
func (r *Requests) RemoveOlderFrom(timestamp time.Time) {
	// start from the beginning and go until timestamp is smaller than the request timestamp

	// index from start to remove
	toRemoveIndex := 0

	for _, val := range (*r).Data {
		if val.Timestamp.Before(timestamp) {
			toRemoveIndex++
		}
	}

	// slice out the ones that are out of the given time
	r.mu.Lock()
	defer r.mu.Unlock()
	(*r).Data = (*r).Data[toRemoveIndex:]
}
