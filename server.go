package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) (int, bool)
	IncrementPlayerScore(name string)
}

type PlayerServer struct {
	store PlayerStore
}

func NewPlayerServer(store PlayerStore) PlayerServer {
	server := PlayerServer{store}
	return server
}

func (s *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	playerName := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodGet:
		s.showScore(w, playerName)
	case http.MethodPost:
		s.processWin(w, playerName)
	}
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
