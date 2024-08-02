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

type Winner struct {
	Name string `json:"name"`
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

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/league/", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.winHandler))
	router.Handle("/add/", http.HandlerFunc(p.addHandler))
	router.Handle("/info/", http.HandlerFunc(p.infoHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) addHandler(w http.ResponseWriter, r *http.Request) {
	var player Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Errorf("problem with json decoding: %v", err)
	}
	p.store.AddPlayer(&player)
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
}

func (p *PlayerServer) winHandler(w http.ResponseWriter, r *http.Request) {
	var winner Winner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&winner); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Errorf("problem with json decoding: %v", err)
	}
	p.processWin(w, winner.Name)
}

func (p *PlayerServer) infoHandler(w http.ResponseWriter, r *http.Request) {
	playerName := strings.TrimPrefix(r.URL.Path, "/info/")
	fmt.Fprint(w, p.store.GetPlayerScore(playerName))
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
