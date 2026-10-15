package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func (s *SpyStore) GetLeague() []Player {
	players := make([]Player, 0)

	for playerName, playerScore := range s.scores {
		players = append(players, Player{playerName, playerScore})
	}

	return players
}

func TestGETPlayers(t *testing.T) {
	store := &SpyStore{scores: map[string]int{
		"George":    6,
		"Anastasia": 7,
	}}
	server := NewPlayerServer(store)

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

func TestGetLeague(t *testing.T) {
	t.Run("league in json format", func(t *testing.T) {
		store := &SpyStore{scores: map[string]int{
			"George":    6,
			"Anastasia": 7,
		}}
		server := NewPlayerServer(store)

		t.Run("check /league endpoint status code and response", func(t *testing.T) {
			request, _ := createNewGetLeagueRequest()
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			want := []Player{{
				Name:  "George",
				Score: 6,
			}, {
				Name:  "Anastasia",
				Score: 7,
			}}

			var got []Player
			err := json.NewDecoder(response.Body).Decode(&got)
			if err != nil {
				t.Errorf("error decoding /league response: %v", err)
			}

			assertStatus(t, response.Code, http.StatusOK)
			assertLeague(t, got, want)
			assertContentType(t, response, "application/json")
		})
	})
}

func createNewGetScoreRequest(name string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	if err == nil {
		return request, nil
	}
	return nil, err
}

func createNewGetLeagueRequest() (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, "/league", nil)
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
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func assertLeague(t *testing.T, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertContentType(t *testing.T, r *httptest.ResponseRecorder, want string) {
	t.Helper()
	if r.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of application/json, got %v", r.Result().Header)
	}
}
