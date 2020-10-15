package main

import (
	// "bytes"
	"crypto/md5"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const jsonContentType = "application/json"

///////////////////////////////////////////////////////////
/// UNIT TESTS

func TestFormModel_GetHash(t *testing.T) {
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

			got := tt.form.GetHash()
			b := md5.Sum([]byte(tt.form.String()))
			want := b2s(b[:])

			// debugTest(t, "got=%s want=%s", got, want)
			assertString(t, got, want)

		})
	}

}

func TestFileSystemFormModel(t *testing.T) {

	tests := []struct {
		name string
		data string
		want FormModel
	}{
		{"Read model", `{"prenom": "patrick", "Nom": "cao"}`, FormModel{Nom: "cao", Prenom: "patrick"}},
		{"Read partial model", `{"Nom": "cao"}`, FormModel{Nom: "cao"}},
		{"Read small model", `{"cao","patrick"}`, FormModel{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datas := strings.NewReader(tt.data)
			store := FileSystemFormModel{datas}

			got, _ := store.Read()
			want := tt.want

			assertFormModel(t, got, want)
		})
	}

	t.Run("MultiRead model", func(t *testing.T) {
		datas := strings.NewReader(`{"prenom": "patrick", "Nom": "cao"}`)
		store := FileSystemFormModel{datas}

		// pass 1
		got, _ := store.Read()
		want := FormModel{Nom: "cao", Prenom: "patrick"}

		assertFormModel(t, got, want)

		// pass 2
		got, _ = store.Read()

		assertFormModel(t, got, want)
	})

	t.Run("Read invalid model", func(t *testing.T) {
		datas := strings.NewReader(`{"Nom": "cao", "Prenom": `)
		store := FileSystemFormModel{datas}

		_, err := store.Read()

		if err == nil {
			t.Errorf("got no error")
		}
	})

}
func TestFormServer_GetResponseAndError(t *testing.T) {
	// get Stub datas
	store := stubTestingStorage{
		map[string]string{
			"/":         "init store",
			"/1/coucou": "email store",
			"/2/coucou": "email validation",
		},
		FormModel{},
		0,
		map[string]FormModel{},
	}
	// debugTest(t,"TestTemplateFromURL/store = %s", store)
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

func TestFormServer_PostResponse(t *testing.T) {
	// get Stub datas
	store := stubTestingStorage{
		map[string]string{
			"/":         "init store",
			"/1/coucou": "email store",
			"/2/coucou": "email validation",
		},
		FormModel{Nom: "cao", Prenom: "patrick"},
		0,
		map[string]FormModel{},
	}
	// debugTest(t,"TestTemplateFromURL/store = %s", store)
	server := NewFormServer(&store)

	tests := []struct {
		name         string
		url          string
		JSONBody     string
		httpCode     int
		returnString string
	}{
		{
			"POST",
			"/",
			`{"nom": "CAO","prenom": "PATRICK"}`,
			http.StatusOK,
			"8ff74764ddedaa020238f31aaf871b22",
		},
		{
			"BAD POST",
			"/",
			`{"CAO","PATRICK`,
			http.StatusBadRequest,
			"",
		},
		{
			"EMPTY POST",
			"/",
			``,
			http.StatusBadRequest,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_"+tt.url, func(t *testing.T) {
			// debugTest(t,"request %s. body='%s'", tt.name, tt.JSONBody)
			request, _ := http.NewRequest(http.MethodPost, tt.url, strings.NewReader(tt.JSONBody))
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)
			// debugTest(t,"response %s. code=%d body=%s", tt.name, response.Code, response.Body.String())

			assertHttpCode(t, response.Code, tt.httpCode)
			assertString(t, response.Body.String(), tt.returnString)
		})
	}
}

func TestFormServer_RecordFormModel(t *testing.T) {
	// get Stub datas
	store := stubTestingStorage{
		map[string]string{
			"/":         "init store",
			"/1/coucou": "email store",
			"/2/coucou": "email validation",
		},
		FormModel{},
		0,
		map[string]FormModel{},
	}
	// debugTest(t,"TestTemplateFromURL/store = %s", store)
	server := NewFormServer(&store)

	tests := []struct {
		name          string
		url           string
		JSONBody      string
		httpCode      int
		spyWriteCalls int
		newModel      FormModel
	}{
		{
			"POST",
			"/",
			`{"nom": "CAO","prenom": "PATRICK"}`,
			http.StatusOK,
			1,
			FormModel{Nom: "CAO", Prenom: "PATRICK"},
		},
		{
			"BAD POST, keep the model",
			"/",
			`{"nom": "CAO","prenom": "PATRI`,
			http.StatusBadRequest,
			1,
			FormModel{Nom: "CAO", Prenom: "PATRICK"},
		},
		{
			"POST, modify the model",
			"/",
			`{"nom": "Escudier","prenom": "Chloe"}`,
			http.StatusOK,
			2,
			FormModel{Nom: "Escudier", Prenom: "Chloe"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_"+tt.url, func(t *testing.T) {
			// debugTest(t,"request %s. body='%s'", tt.name, tt.JSONBody)
			request, _ := http.NewRequest(http.MethodPost, tt.url, strings.NewReader(tt.JSONBody))
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			// debugTest(t,"response %s. code=%d body=%s", tt.name, response.Code, response.Body.String())

			assertHttpCode(t, response.Code, tt.httpCode)
			if tt.spyWriteCalls != store.spyWriteCalls {
				t.Errorf("got %d calls to RecordFormModel want %d", store.spyWriteCalls, tt.spyWriteCalls)
			}
			assertFormModel(t, store.GetModel(), tt.newModel)
			// assertString(t, response.Body.String(), tt.returnString)
		})
	}
}

///////////////////////////////////////////////////////////
/// INTEGRATION TESTS

func TestFormServer_Integration(t *testing.T) {
	t.Skip("TODO")
	// get Stub datas
	store := stubTestingStorage{
		map[string]string{},
		FormModel{},
		0,
		map[string]FormModel{},
	}
	// debugTest(t,"TestTemplateFromURL/store = %s", store)
	server := NewFormServer(&store)

	postURL := "/"
	JSONBody := `{"nom": "cao","prenom": "patrick"}`

	getURL := "/1/"
	// wantModel := FormModel{Nom: "CAO", Prenom: "PATRICK"}
	httpCode := http.StatusOK

	t.Run("Record and get model", func(t *testing.T) {
		// REQ1 : POST with body
		debugTest(t, "request  1. body='%s'", JSONBody)
		request, _ := http.NewRequest(http.MethodPost, postURL, strings.NewReader(JSONBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		debugTest(t, "response 1. code=%d body=%s", response.Code, response.Body.String())

		hash := response.Body.String()

		getURL += hash
		// getURL = url.PathEscape(getURL)

		// REQ2 : GET with hash in URL
		debugTest(t, "request  2. url='%v'", getURL)
		request, _ = http.NewRequest(http.MethodGet, getURL, nil)
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)
		debugTest(t, "response 2. code=%d body=%s", response.Code, response.Body.String())

		assertHttpCode(t, response.Code, httpCode)
		model, err := ReadFormModel(response.Body)
		if err != nil {
			t.Errorf("error when parsing JSON %v", err)
		}
		assertFormModel(t, model, FormModel{Nom: "cao", Prenom: "patrick"})
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

func assertBytes(t *testing.T, got, want []byte) {
	t.Helper()

	if string(got) != string(want) {
		t.Errorf("got %x want %x", got, want)
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

var assertContentTypeJSON = func(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	assertContentType(t, response, jsonContentType)
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

///////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////

// debug call log.Printf with prefix DEBUG
func debug(f string, i ...interface{}) {
	prefix := fmt.Sprintf("DEBUG %s", f)
	log.Printf(prefix, i...)
}
func debugTest(t *testing.T, f string, i ...interface{}) {
	prefix := fmt.Sprintf("[%s] %v", t.Name(), f)
	debug(prefix, i...)
}

///////////////////////////////////////////////////////////
// stubTestingStorage
type stubTestingStorage struct {
	html          map[string]string
	model         FormModel
	spyWriteCalls int
	hashs         map[string]FormModel
}

func (s *stubTestingStorage) TemplateFromURL(aUrl string) string {
	if someHtml, err := s.html[aUrl]; err {
		// debug(t,">> TemplateFromURL(%s) in html = %s", aUrl, someHtml)
		return someHtml
	} else {
		// debug(">> TemplateFromURL(%s) not in html", aUrl)
		return ""
	}
}
func (s *stubTestingStorage) GetModel() FormModel {
	// debug(">> GetModel()")
	return s.model
}

func (s *stubTestingStorage) RecordModel(model FormModel) error {
	s.spyWriteCalls++
	s.model = model
	return nil
}

func (s *stubTestingStorage) GetHashs() map[string]FormModel {
	// debug(">> GetHashs()")
	return s.hashs
}
