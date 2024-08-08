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

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[int]int{
			1: 20,
			2: 10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "The player with id: 1 has 20 wins")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest(2)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "The player with id: 2 has 10 wins")
	})

	t.Run("returns 0 on missing players", func(t *testing.T) {
		request := newGetScoreRequest(3)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "The player with id: 3 has 0 wins")
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		make(map[int]int),
		make([]int, 0),
		make([]Player, 0),
	}
	server := NewPlayerServer(&store)

	t.Run("it records wins on POST", func(t *testing.T) {
		player := 1

		server.ServeHTTP(httptest.NewRecorder(), newPlayerCreateRequest(player, "Test", 3))

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		AssertPlayerWin(t, &store, player)
	})
}

func TestLeague(t *testing.T) {

	t.Run("it returns the League table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{1, "Cleo", 32},
			{2, "Chris", 20},
			{3, "Tiest", 14},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)

	})
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
