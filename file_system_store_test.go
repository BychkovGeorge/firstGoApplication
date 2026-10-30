package main

import (
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	file, removeFileFunc := assertCreateFile(t, `[
			{"Name": "Cleo", "Score": 10},
			{"Name": "Chris", "Score": 33}]`)
	store, err := NewFileSystemPlayerStore(file)
	assertNoError(t, err)
	defer removeFileFunc()

	t.Run("league from a reader", func(t *testing.T) {
		got := store.GetLeague()

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("score from a reader", func(t *testing.T) {
		got := store.GetPlayerScore("Cleo")
		want := 10

		assertScoreEquals(t, got, want)
	})

	t.Run("league sorted", func(t *testing.T) {
		got := store.GetLeague()

		want := League{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("record new win", func(t *testing.T) {
		store.IncrementPlayerScore("Cleo")
		got := store.GetPlayerScore("Cleo")
		want := 11

		assertScoreEquals(t, got, want)
	})

	t.Run("record win for new player", func(t *testing.T) {
		store.IncrementPlayerScore("Anastasia")
		got := store.GetPlayerScore("Anastasia")
		want := 1

		assertScoreEquals(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, removeFile := assertCreateFile(t, "")
		defer removeFile()

		_, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)
	})
}

func assertCreateFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	file, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	file.Write([]byte(initialData))

	removeFile := func() {
		file.Close()
		os.Remove(file.Name())
	}

	return file, removeFile
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("error was not expecting any errors: %v", err)
	}
}
