package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	// "log"
	"net/http"
	// "strings"
)

// Data Model
type FormModel struct {
	Nom    string
	Prenom string
}

func (f *FormModel) String() string {
	return fmt.Sprintf("Nom: %s, Prenom: %s", f.Nom, f.Prenom)
}

func (f *FormModel) getHash() string {
	h := md5.New()
	io.WriteString(h, f.String())
	return string(h.Sum(nil))
}

// FIXME
// func (f *FormModel) String() string {
// 	s, err := json.Marshal(f)
// 	if err != nil {
// 		return ""
// 	}
// 	return string(s)
// }

// struct FormServer
// is a http.Handler !!
type FormServer struct {
	form GetterForm
	http.Handler
}

type GetterForm interface {
	// get html from url
	GetForm(string) string
	GetModel() FormModel
}

// Constructor
func NewFormServer(form GetterForm) *FormServer {

	f := new(FormServer)
	f.form = form

	router := http.NewServeMux()

	router.Handle("/", http.HandlerFunc(f.handle_step0))
	router.Handle("/forms", http.HandlerFunc(f.handle_json_forms))

	f.Handler = router

	return f
}

//
// Handlers

// default
func (f *FormServer) handle_step0(w http.ResponseWriter, r *http.Request) {
	url := r.URL.RequestURI()
	resp := f.form.GetForm(url)
	if resp == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, resp)
}

// API - json
// /forms
func (f *FormServer) handle_json_forms(w http.ResponseWriter, r *http.Request) {
	// log.Printf(">> handle_json_forms")
	w.Header().Set("content-type", "application/json")
	model := f.form.GetModel()
	// log.Printf("   handle_json_forms/GetModel = %s", model)
	json.NewEncoder(w).Encode(model)
	// http: superfluous
	// w.WriteHeader(http.StatusOK)
	// log.Printf("<< handle_json_forms")
}
