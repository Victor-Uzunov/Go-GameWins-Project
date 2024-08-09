package poker

import (
	"io"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &tape{file}

	if _, err := tape.Write([]byte("abc")); err != nil {
		t.Fatal(err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
