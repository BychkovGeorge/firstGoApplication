package main

import (
	"log"
	"net/http"
	"sync"
)

type Player struct {
	Name  string
	Score int
}

type InMemoryPlayerStore struct {
	scores map[string]int
	mu     sync.Mutex
}

func (s *InMemoryPlayerStore) GetPlayerScore(name string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.scores == nil {
		s.scores = make(map[string]int)
	}

	score, ok := s.scores[name]
	if ok {
		return score, true
	}
	return 0, false
}

func (s *InMemoryPlayerStore) IncrementPlayerScore(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.scores == nil {
		s.scores = make(map[string]int)
	}

	score, ok := s.scores[name]
	if ok {
		s.scores[name] = score + 1
	} else {
		s.scores[name] = 1
	}
}

func (s *InMemoryPlayerStore) GetLeague() []Player {
	s.mu.Lock()
	defer s.mu.Unlock()

	players := make([]Player, 0)

	for playerName, playerScore := range s.scores {
		players = append(players, Player{playerName, playerScore})
	}

	return players
}

func main() {
	playerServer := NewPlayerServer(&InMemoryPlayerStore{
		scores: map[string]int{
			"Anastasia": 10,
			"George":    2,
		},
		mu: sync.Mutex{},
	})
	handler := http.HandlerFunc(playerServer.ServeHTTP)
	log.Fatal(http.ListenAndServe(":5001", handler))
}
