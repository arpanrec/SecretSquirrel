package main

import (
	"gitlab.com/arpanrecme/initsecureserver/internal/iss"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		iss.EntryPoint(w, r)

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
