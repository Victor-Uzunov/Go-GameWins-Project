package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string) error
	GetLeague() League
	AddPlayer(player *Player) error
	DeletePlayer(name string) error
}

type Winner struct {
	Name string `json:"name"`
}

type Resource struct {
	Name      string `json:"name"`
	Wins      int    `json:"wins"`
	CreatedAt string `json:"created_at"`
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
	router.Handle("/update/", http.HandlerFunc(p.updateHandler))
	router.Handle("/create/", http.HandlerFunc(p.createHandler))
	router.Handle("/info/", http.HandlerFunc(p.infoHandler))
	router.Handle("/delete/", http.HandlerFunc(p.deleteHandler))

	p.Handler = router

	return p
}

// DELETE
func (p *PlayerServer) deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	username := parts[2]
	if err := p.store.DeletePlayer(username); err != nil {
		http.Error(w, "Invalid username provided", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// POST
func (p *PlayerServer) createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var player Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := p.store.AddPlayer(&player); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	resource := Resource{
		Name:      player.Name,
		Wins:      player.Wins,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	w.Header().Set("Location", "/info/"+player.Name)
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resource); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// GET
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("content-type", jsonContentType)
	if err := json.NewEncoder(w).Encode(p.store.GetLeague()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PATCH
func (p *PlayerServer) updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var winner Winner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&winner); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	p.processWin(w, winner.Name)
}

// GET
func (p *PlayerServer) infoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	playerName := strings.TrimPrefix(r.URL.Path, "/info/")
	if _, err := fmt.Fprintf(w, "The player %s has %d wins", playerName, p.store.GetPlayerScore(playerName)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	if err := p.store.RecordWin(player); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, "The player with name %s has %d wins now", player, p.store.GetPlayerScore(player)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
