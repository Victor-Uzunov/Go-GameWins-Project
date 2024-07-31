package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
	AddPlayer(player *Player)
}

type Player struct {
	Name string `json:"name"`
	Wins int    `json:"wins"`
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

const jsonContentType = "application/json"

// NewPlayerServer creates a PlayerServer with routing configured.
func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/league/", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/add/", http.HandlerFunc(p.addHandler))
	router.Handle("/info/", http.HandlerFunc(p.infoHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) addHandler(w http.ResponseWriter, r *http.Request) {
	var player Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	p.store.AddPlayer(&player)
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	var player Player
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&player); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	p.processWin(w, player.Name)
}

func (p *PlayerServer) infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	playerName := strings.TrimPrefix(r.URL.Path, "/info/")
	json.NewEncoder(w).Encode(p.store.GetLeague().Find(playerName))
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusOK)
}
