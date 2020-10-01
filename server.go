package main

import (
	// "bytes"
	"log"
	"net/http"
)

var (
// buf   bytes.Buffer
// mylog = log.New(&buf, "inscription: ", log.Lshortfile)
)

func main() {
	// TODO get port from command line
	port := "5000"

	handler := http.HandlerFunc(FormServer)
	log.Printf("Running server on port %d ...\n", 5000)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
