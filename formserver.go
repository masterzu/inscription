package main

import (
	"fmt"
	// "log"
	"net/http"
	"strings"
)

// steps for process validation
// https://yourbasic.org/golang/iota/
type step int

const (
	STEP0 step = iota
	STEP1
	STEP2
	STEP_ERROR
)

func (s step) String() string {
	return [...]string{"Step 0", "Step 1", "Step 2", "Step unknown"}[s]
}

func FormServer(w http.ResponseWriter, request *http.Request) {

	url := request.URL.RequestURI()
	fmt.Fprint(w, getForm(url))

}

func getForm(url string) string {
	switch getFormStep(url) {

	case STEP0:
		return "init form"
	case STEP1:
		return "email form"
	case STEP2:
		return "valid email form"
	default:
		return "42"
	}

}

func getFormStep(url string) step {
	if url == "/" {
		return STEP0
	}
	if strings.HasPrefix(url, "/1/") {
		return STEP1
	} else if strings.HasPrefix(url, "/2/") {
		return STEP2
	}
	return STEP_ERROR

}
