package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := InMemoryPlayerStore{}
	server := NewPlayerServer(&store)
	player := "Pepper"

	getRequest, _ := createNewGetScoreRequest(player)
	postRequest, _ := createNewPostScoreRequest(player)

	var wg = sync.WaitGroup{}
	wg.Add(1000)
	for range 1000 {
		go func() {
			server.ServeHTTP(httptest.NewRecorder(), postRequest)
			wg.Done()
		}()
	}
	wg.Wait()

	response := httptest.NewRecorder()
	server.ServeHTTP(response, getRequest)

	assertStatus(t, response.Code, http.StatusOK)
	assertRequestBody(t, response.Body.String(), "1000")
}
