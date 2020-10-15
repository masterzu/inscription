package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// = steps for process validation
type FormServer struct {
	store Storage
	http.Handler
}
// + https://yourbasic.org/golang/iota/
// const (
func NewFormServer(st Storage) *FormServer {
// 	STEP0 step = iota
	f := new(FormServer)
	f.store = st
	router := http.NewServeMux()
// 	STEP1
	router.Handle("/", http.HandlerFunc(f.handleStep0))
	router.Handle("/1/", http.HandlerFunc(f.handleStep1))
// 	STEP2
// 	STEP_ERROR
	router.Handle("/forms", http.HandlerFunc(f.handleJSONForms))
// )
	f.Handler = router
// func (s step) String() string {
	return f
}
// 	return [...]string{"Step 0", "Step 1", "Step 2", "Step unknown"}[s]
// }
// func getFormStep(url string) step {
// 	if url == "/" {
// 		return STEP0
func (f *FormServer) handleStep0(w http.ResponseWriter, r *http.Request) {
// 	}
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

// func write(url string) string {
func (f *FormServer) handleStep1(w http.ResponseWriter, r *http.Request) {
// 	switch getFormStep(url) {

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

