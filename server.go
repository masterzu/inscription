package main

import (
	"log"
	"net/http"
)

func main() {
	// TODO get port from command line
	port := "5000"

	server := NewFormServer(NewMemoryStore())

	// TODO print access log
	log.Printf("Running server on port %d ...\n", 5000)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
