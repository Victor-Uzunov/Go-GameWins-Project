package poker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	err := initialisePlayerDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
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

func (f *FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(id int) int {

	player := f.league.Find(id)

	if player != nil {
		return player.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWin(id int) error {

	player := f.league.Find(id)

	if player != nil {
		player.Wins++
	} else {
		return errors.New("player with this id not found")
	}

	return f.database.Encode(f.league)
}

func (f *FileSystemPlayerStore) AddPlayer(player *Player) error {
	if player.Name == "" {
		return errors.New("player name cannot be empty")
	}
	f.league = append(f.league, *player)
	if err := f.database.Encode(f.league); err != nil {
		return errors.New("failed to add player to database" + err.Error())
	}
	return nil
}

func (f *FileSystemPlayerStore) DeletePlayer(id int) error {
	for i, player := range f.league {
		if player.ID == id {
			f.league = removeElement(f.league, i)
			if err := f.database.Encode(f.league); err != nil {
				return errors.New("failed to remove player from database" + err.Error())
			}
			return nil
		}
	}
	return errors.New("player not found")
}

func removeElement(slice []Player, index int) []Player {
	return append(slice[:index], slice[index+1:]...)
}
