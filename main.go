package main

import (
	"log"
	"net/http"
)

func main() {

	// Serve all files in the static folder; automatically serves index.html at "/"
	http.Handle("/", http.FileServer(http.Dir("static")))

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}