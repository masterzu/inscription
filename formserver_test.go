package main

import (
	"crypto/md5"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	// "log"
)

const jsonContentType = "application/json"

///////////////////////////////////////////////////////////
/// TESTS

func TestGetHash(t *testing.T) {
	tests := []struct {
		name string
		form FormModel
	}{
		{"patrick cao", FormModel{Nom: "cao", Prenom: "patrick"}},
		{"patrick", FormModel{Prenom: "patrick"}},
		{"", FormModel{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := tt.form.getHash()
			h := md5.New()
			s := tt.form.String()
			io.WriteString(h, s)
			want := string(h.Sum(nil))

			assertString(t, got, want)

		})
	}

}

func TestGetForm(t *testing.T) {
	// get Stub datas
	store := stubGetterForm{
		map[string]string{
			"/":         "init store",
			"/1/coucou": "email store",
			"/2/coucou": "email validation",
		},
		FormModel{},
	}
	// log.Printf("TestGetForm/store = %s", store)
	server := NewFormServer(&store)

	tests := []struct {
		name         string
		url          string
		httpCode     int
		returnString string
	}{
		{name: "root URL", url: "/", httpCode: http.StatusOK, returnString: "init store"},
		{name: "URL step1", url: "/1/coucou", httpCode: http.StatusOK, returnString: "email store"},
		{name: "URL step2", url: "/2/coucou", httpCode: http.StatusOK, returnString: "email validation"},
		{name: "URL unknown", url: "/1/qwe", httpCode: http.StatusNotFound, returnString: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertString(t, response.Body.String(), tt.returnString)
			assertHttpCode(t, response.Code, tt.httpCode)
		})

	}

}

func TestPostForm(t *testing.T) {
	t.Skip()
	// get Stub datas
	store := stubGetterForm{
		map[string]string{
			"/":         "init store",
			"/1/coucou": "email store",
			"/2/coucou": "email validation",
		},
		FormModel{},
	}
	// log.Printf("TestGetForm/store = %s", store)
	server := NewFormServer(&store)

	url := "/1/"

	t.Run("/post", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, url, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		// assertString(t, response.Body.String(), tt.returnString)
		assertHttpCode(t, response.Code, http.StatusOK)
	})
}

func TestGetJsonForms(t *testing.T) {
	// t.Skip("todo")
	// get Stub datas

	model := FormModel{
		Nom:    "Cao",
		Prenom: "Patrick",
	}
	// log.Printf("TestGetJsonForms/model = %s", model)
	store := stubGetterForm{
		map[string]string{},
		model,
	}
	server := NewFormServer(&store)

	t.Run("all datas", func(t *testing.T) {

		url := "/forms"

		request, _ := http.NewRequest(http.MethodGet, url, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if response.Result().Header.Get("content-type") != "application/json" {
			t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
		}

		assertHttpCode(t, response.Code, http.StatusOK)

		got := getFormModelFromResponse(t, response.Body)
		assertFormModel(t, got, model)
	})
}

///////////////////////////////////////////////////////////
/// Asserts
/// Helpers

func assertString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertInt(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

var assertHttpCode = assertInt

func assertFormModel(t *testing.T, got, want FormModel) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

///////////////////////////////////////////////////////////
/// Helpers

func getFormModelFromResponse(t *testing.T, body io.Reader) (model FormModel) {
	var got FormModel
	err := json.NewDecoder(body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into []FormModel, '%v'", body, err)
	}
	return got
}

// stubGetterForm
type stubGetterForm struct {
	html  map[string]string
	model FormModel
}

func (s *stubGetterForm) GetForm(url string) string {
	if someHtml, err := s.html[url]; err {
		// log.Printf(">> GetForm(%s) in html = %s", url, someHtml)
		return someHtml
	} else {
		// log.Printf(">> GetForm(%s) not in html", url)
		return ""
	}
}
func (s *stubGetterForm) GetModel() FormModel {
	// log.Printf(">> GetModel()")
	return s.model
}
