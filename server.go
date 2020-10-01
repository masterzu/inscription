package main

import (
	// "bytes"
	"log"
	"net/http"
)

type stubForm struct{}

func (s *stubForm) write(url string) string {
	return "all is LOVE"
}

func main() {
	// TODO get port from command line
	port := "5000"

	server := &FormServer{&stubForm{}}

	log.Printf("Running server on port %d ...\n", 5000)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
