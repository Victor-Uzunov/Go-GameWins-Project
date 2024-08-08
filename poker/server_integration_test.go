package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := NewFileSystemPlayerStore(database)

	assertNoError(t, err)

	server := NewPlayerServer(store)
	id := 0
	server.ServeHTTP(httptest.NewRecorder(), newPlayerCreateRequest(id, "Test", 3))

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(id))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(id))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(id))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(id))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "The player with id: 0 has 6 wins")
	})

	t.Run("get League", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{0, "Test", 6},
		}
		assertLeague(t, got, want)
	})
}
