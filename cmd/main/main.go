package main

import (
	"log"
	"net/http"

	"gitlab.com/arpanrecme/initsecureserver/internal/iss"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		iss.EntryPoint(w, r)

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
