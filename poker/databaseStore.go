package poker

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

type DatabaseStore struct {
	db *sqlx.DB
}

func NewDatabaseStore(conStr string) (*DatabaseStore, error) {
	db, err := sqlx.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	store := &DatabaseStore{db: db}
	return store, nil
}

func (store *DatabaseStore) GetLeague() League {
	var league League
	err := store.db.Select(&league, "SELECT p.username as name, COUNT(gr.winner_id) AS wins\nFROM players AS p\nJOIN game_results AS gr ON p.id = gr.winner_id\nGROUP BY p.username\nORDER BY wins DESC")
	if err != nil {
		return nil
	}
	return league
}

func (store *DatabaseStore) GetPlayerScore(name string) int {
	var wins int
	err := store.db.Get(&wins, "SELECT COUNT(gr.winner_id) AS wins\nFROM players AS p\n         JOIN game_results AS gr ON p.id = gr.winner_id\nWHERE p.username = $1\nGROUP BY p.username\nORDER BY wins DESC;", name)
	if err != nil {
		return 0
	}
	return wins
}

func (store *DatabaseStore) RecordWin(name string) error {
	if name == "" {
		return errors.New("no name provided")
	}
	return nil
	//TODO
}
func (store *DatabaseStore) AddPlayer(player *Player) error {
	//TODO
	return nil
}

func (store *DatabaseStore) DeletePlayer(name string) error {
	//TODO
	return nil
}
