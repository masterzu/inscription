package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func gotWantString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func gotWantInt(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

// stubWriterString
// to store links url -> html
type stubWriterString struct {
	pages map[string]string
}

func (s *stubWriterString) write(url string) string {
	return s.pages[url]
}

///////////////////////////////////////////////////////////
// TESTS

func TestGetForm(t *testing.T) {
	form := stubWriterString{
		map[string]string{
			"/":         "init form",
			"/1/coucou": "email form",
			"/2/toto":   "email validation",
		},
	}

	server := &FormServer{&form}

	urlWant := func(t *testing.T, url, want string) {
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s", url), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		gotWantString(t, got, want)
	}

	t.Run("empty URL", func(t *testing.T) {
		urlWant(t, "/", "init form")
	})

	t.Run("/1/ URL", func(t *testing.T) {
		urlWant(t, "/1/coucou", "email form")
	})

	t.Run("/2/ URL", func(t *testing.T) {
		urlWant(t, "/2/toto", "email validation")
	})

}
