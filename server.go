package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) (int, bool)
	IncrementPlayerScore(name string)
	GetLeague() []Player
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	server := new(PlayerServer)
	server.store = store

	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(server.handleLeague))
	router.Handle("/players/", http.HandlerFunc(server.handlePlayers))

	server.Handler = router

	return server
}

func (s *PlayerServer) showScore(w http.ResponseWriter, name string) {
	score, ok := s.store.GetPlayerScore(name)
	if ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}

func (s *PlayerServer) processWin(w http.ResponseWriter, name string) {
	s.store.IncrementPlayerScore(name)
	w.WriteHeader(http.StatusAccepted)
}

func (s *PlayerServer) handleLeague(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.store.GetLeague())
}

func (s *PlayerServer) handlePlayers(w http.ResponseWriter, r *http.Request) {
	playerName := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodGet:
		s.showScore(w, playerName)
	case http.MethodPost:
		s.processWin(w, playerName)
	}
}
