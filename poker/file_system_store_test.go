package poker

import (
	"os"
	"testing"
)

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	if _, err = tmpFile.Write([]byte(initialData)); err != nil {
		t.Fatalf("could not write to temp file %v", err)
	}

	removeFile := func() {
		if err = tmpFile.Close(); err != nil {
			t.Fatalf("could not close temp file %v", err)
		}
		if err = os.Remove(tmpFile.Name()); err != nil {
			t.Fatalf("could not remove temp file %v", err)
		}
	}

	return tmpFile, removeFile
}

func TestFileSystemStore(t *testing.T) {

	t.Run("League sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"id": 2, "name": "Cleo", "wins": 10},
			{"id": 1, "name": "Chris", "wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetLeague()

		want := []Player{
			{1, "Chris", 33},
			{2, "Cleo", 10},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"ID": 1, "Name": "Cleo", "Wins": 10},
			{"ID": 2, "Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetPlayerScore(2)
		want := 33
		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"ID": 1, "Name": "Cleo", "Wins": 10},
			{"ID": 2, "Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		store.RecordWin(2)

		got := store.GetPlayerScore(2)
		want := 34
		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"ID": 1, "Name": "Cleo", "Wins": 10},
			{"ID": 2, "Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		err = store.RecordWin(3)
		assertError(t, err)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)
	})
}

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

func assertError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}
}
