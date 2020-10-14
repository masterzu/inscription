package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FormServer is the HTTP.Handler for the app
type FormServer struct {
	store Storage
	http.Handler
}

// Constructor
func NewFormServer(st Storage) *FormServer {

	f := new(FormServer)
	f.store = st
	router := http.NewServeMux()

	router.Handle("/", http.HandlerFunc(f.handleStep0))
	router.Handle("/1/", http.HandlerFunc(f.handleStep1))

	//api
	router.Handle("/forms", http.HandlerFunc(f.handleJSONForms))

	f.Handler = router

	return f
}

//
// Handlers

// URL = /
func (f *FormServer) handleStep0(w http.ResponseWriter, r *http.Request) {
	// debug(">> handleStep0/method: %s", r.Method)
	aUrl := r.URL.RequestURI()
	switch r.Method {
	case "GET":
		resp := f.store.TemplateFromURL(aUrl)
		if resp == "" {
			w.WriteHeader(http.StatusNotFound)
		}
		fmt.Fprint(w, resp)
	case "POST":
		model, err := ReadFormModel(r.Body)
		if err != nil {
			// debug("handleStep0/POST/err=%s", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "")
			return
		}
		hash := model.GetHash()
		f.store.RecordModel(model, string(hash))
		fmt.Fprint(w, string(hash))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "")
	}
}

// URL = /1/
func (f *FormServer) handleStep1(w http.ResponseWriter, r *http.Request) {
	// debug(">> handleStep1/method: %s", r.Method)

	switch r.Method {
	case "GET":
		aUrl := r.URL.RequestURI()
		resp := f.store.TemplateFromURL(aUrl)
		if resp == "" {
			w.WriteHeader(http.StatusNotFound)
		}
		fmt.Fprint(w, resp)

	case "POST":

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "")
	}
}

// API - json
// GET /forms
func (f *FormServer) handleJSONForms(w http.ResponseWriter, r *http.Request) {
	// debug(">> handleJSONForms")
	w.Header().Set("content-type", "application/json")
	model := f.store.GetModel()
	// debug("   handleJSONForms/GetModel = %s", model)
	json.NewEncoder(w).Encode(model)
	// http: superfluous
	// w.WriteHeader(http.StatusOK)
	// debug("<< handleJSONForms")
}
