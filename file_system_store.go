package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
	mutex    sync.Mutex
}

func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	store, err := NewFileSystemPlayerStore(db)

	if err != nil {
		return nil, nil, fmt.Errorf("problem creating file system player store, %v ", err)
	}

	return store, closeFunc, nil
}

func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, io.SeekStart)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}

	return nil
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	err := initialisePlayerDBFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem initialising player DB file %s, %v", file.Name(), err)
	}

	league, err := NewLeague(file)
	if err != nil {
		return nil, fmt.Errorf("could not create league: %v", err)
	}
	var mutex sync.Mutex
	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
		mutex:    mutex,
	}, nil
}

func (store *FileSystemPlayerStore) GetLeague() League {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	sort.Slice(store.league, func(i, j int) bool {
		return store.league[i].Score > store.league[j].Score
	})
	return store.league
}

func (store *FileSystemPlayerStore) GetPlayerScore(name string) int {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	player := store.league.Find(name)
	if player != nil {
		return player.Score
	}
	return 0
}

func (store *FileSystemPlayerStore) IncrementPlayerScore(name string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	player := store.league.Find(name)
	if player != nil {
		player.Score++
	} else {
		store.league = append(store.league, Player{name, 1})
	}
	store.database.Encode(store.league)
}
