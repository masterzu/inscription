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

func TestGetForm(t *testing.T) {

	urlWant := func(t *testing.T, url, want string) {
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s", url), nil)
		response := httptest.NewRecorder()

		FormServer(response, request)

		got := response.Body.String()
		gotWantString(t, got, want)
	}

	t.Run("empty URL", func(t *testing.T) {
		urlWant(t, "/", "init form")
	})

	t.Run("/ URL", func(t *testing.T) {
		urlWant(t, "/", "init form")
	})

	t.Run("/1/ URL", func(t *testing.T) {
		urlWant(t, "/1/qwerq", "email form")
	})

}

func TestGetFormStep(t *testing.T) {

	gotWant := func(t *testing.T, got, want step) {
		gotWantInt(t, int(got), int(want))
	}

	t.Run("getFormStep /", func(t *testing.T) {
		gotWant(t, getFormStep("/"), STEP0)
	})

}
