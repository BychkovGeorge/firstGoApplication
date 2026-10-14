package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SpyStore struct {
	scores map[string]int
}

func (s *SpyStore) GetPlayerScore(name string) (int, bool) {
	score, ok := s.scores[name]
	if ok {
		return score, true
	}
	return 0, false
}

func (s *SpyStore) IncrementPlayerScore(name string) {
	score, ok := s.scores[name]
	if ok {
		s.scores[name] = score + 1
	}
	s.scores[name] = 1
}

func TestGETPlayers(t *testing.T) {
	store := SpyStore{scores: map[string]int{
		"George":    6,
		"Anastasia": 7,
	}}
	server := &PlayerServer{&store}

	t.Run("return George's score", func(t *testing.T) {
		request, _ := createNewGetScoreRequest("George")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertRequestBody(t, response.Body.String(), "6")
	})

	t.Run("return Anastasia's score", func(t *testing.T) {
		request, _ := createNewGetScoreRequest("Anastasia")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertRequestBody(t, response.Body.String(), "7")
	})

	t.Run("return unexisting player's score", func(t *testing.T) {
		request, _ := createNewGetScoreRequest("Peter")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func createNewGetScoreRequest(name string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	if err == nil {
		return request, nil
	}
	return nil, err
}

func createNewPostScoreRequest(name string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	if err == nil {
		return request, nil
	}
	return nil, err
}

func assertRequestBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
