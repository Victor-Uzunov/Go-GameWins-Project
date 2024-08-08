package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

type League []Player

func (l League) Find(id int) *Player {
	for i, p := range l {
		if p.ID == id {
			return &l[i]
		}
	}
	return nil
}

func NewLeague(rdr io.Reader) (League, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)

	if err != nil {
		err = fmt.Errorf("problem parsing League, %v", err)
	}

	return league, err
}
