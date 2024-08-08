package poker

import (
	"fmt"
	"testing"
	"time"
)

type StubPlayerStore struct {
	Scores   map[int]int
	WinCalls []int
	League   []Player
}

func (s *StubPlayerStore) Find(id int) bool {
	_, ok := s.Scores[id]
	return ok
}

func (s *StubPlayerStore) DeletePlayer(id int) error {
	delete(s.Scores, id)
	return nil
}

func (s *StubPlayerStore) GetPlayerScore(id int) int {
	score := s.Scores[id]
	return score
}

func (s *StubPlayerStore) RecordWin(id int) error {
	s.WinCalls = append(s.WinCalls, id)
	return nil
}

func (s *StubPlayerStore) AddPlayer(player *Player) error {
	s.Scores[player.ID] = player.Wins
	return nil
}

func (s *StubPlayerStore) GetLeague() League {
	return s.League
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winnerID int) {
	t.Helper()

	if len(store.WinCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.WinCalls), 1)
	}

	if store.WinCalls[0] != winnerID {
		t.Errorf("did not store correct winner got %q want %q", store.WinCalls[0], winnerID)
	}
}

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

type SpyBlindAlerter struct {
	Alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.Alerts = append(s.Alerts, ScheduledAlert{at, amount})
}
