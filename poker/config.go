package poker

import "time"

type PlayerConfig struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type GameConfig struct {
	ID       int       `db:"id" json:"id"`
	GameDate time.Time `db:"game_date" json:"game_date"`
	Location string    `db:"location" json:"location"`
	Notes    string    `db:"notes" json:"notes"`
}

type ResultConfig struct {
	ID        int `db:"id" json:"id"`
	GameID    int `db:"game_id" json:"game_id"`
	WinnerID  int `db:"winner_id" json:"winner_id"`
	AmountWon int `db:"amount_won" json:"amount_won"`
}
