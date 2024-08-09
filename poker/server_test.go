package poker

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestGetPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[int]int{
			1: 20,
			2: 10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)
	tests := []struct {
		name           string
		playerID       int
		expectedStatus int
		expectedBody   string
	}{
		{"returns player with ID:1 score", 1, http.StatusOK, "The player with id: 1 has 20 wins"},
		{"returns player with ID:2 score", 2, http.StatusOK, "The player with id: 2 has 10 wins"},
		{"returns 0 on missing players", 3, http.StatusOK, "The player with id: 3 has 0 wins"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := newGetScoreRequest(tt.playerID)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)
			assertStatus(t, response.Code, tt.expectedStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedBody)
		})
	}
}

func TestScoreWins(t *testing.T) {
	store := StubPlayerStore{
		make(map[int]int),
		make([]int, 0),
		make([]Player, 0),
	}
	server := NewPlayerServer(&store)
	tests := []struct {
		name           string
		playerID       int
		playerName     string
		wins           int
		expectedStatus int
	}{
		{"records wins on POST", 1, "Test", 3, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := newPlayerCreateRequest(tt.playerID, tt.playerName, tt.wins)
			server.ServeHTTP(httptest.NewRecorder(), request)

			request = newPostWinRequest(tt.playerID)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			assertStatus(t, response.Code, tt.expectedStatus)
			AssertPlayerWin(t, &store, tt.playerID)
		})
	}
}

func TestLeague(t *testing.T) {
	tests := []struct {
		name           string
		league         []Player
		expectedStatus int
		expectedLeague []Player
	}{
		{"returns the league table as JSON", []Player{
			{1, "Test1", 32},
			{2, "Test2", 20},
			{3, "Test3", 14},
		}, http.StatusOK, []Player{
			{1, "Test1", 32},
			{2, "Test2", 20},
			{3, "Test3", 14}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := StubPlayerStore{nil, nil, tt.league}
			server := NewPlayerServer(&store)
			request := newLeagueRequest()
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)
			got := getLeagueFromResponse(t, response.Body)

			assertStatus(t, response.Code, tt.expectedStatus)
			assertLeague(t, got, tt.expectedLeague)
			assertContentType(t, response, jsonContentType)
		})
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Header().Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func getLeagueFromResponse(t testing.TB, body io.Reader) []Player {
	t.Helper()
	league, err := NewLeague(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return league
}

func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league/", nil)
	return req
}

func newGetScoreRequest(id int) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/info/%d", id), nil)
	return req
}

func newPostWinRequest(id int) *http.Request {
	body := fmt.Sprintf(`{"id": %v, "name": "Test"}`, id)
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/update/"), strings.NewReader(body))
	return req
}

func newPlayerCreateRequest(id int, name string, wins int) *http.Request {
	body := fmt.Sprintf(`{"id": %v, "name": "%s", "wins": %v}`, id, name, wins)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/create/"), strings.NewReader(body))
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
