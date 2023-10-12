package main

import (
	"gitlab.com/arpanrecme/initsecureserver/internal/ops"
	"io"
	"log"
	"net/http"
	"strings"
)

func EntryPoint(w http.ResponseWriter, r *http.Request) {

	urlPath := r.URL.Path

	body, errReadAll := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		errBodyClose := Body.Close()
		if errBodyClose != nil {
			log.Fatal(errBodyClose)
		}
	}(r.Body)
	if errReadAll != nil {
		log.Fatal(errReadAll)
	}

	rMethod := r.Method

	query := r.URL.Query()

	header := r.Header

	formData := r.Form

	log.Println("URL Path: ", urlPath, "\nMethod: ", rMethod, "\nHeader: ", header,
		"\nForm Data: ", formData,
		"\nBody: ", string(body), "\nQuery: ", query)

	if strings.HasPrefix(urlPath, "/tfstate/") {
		ops.TerraformStateHandler(string(body), rMethod, urlPath, query, w)
	} else if strings.HasPrefix(urlPath, "/files/") {
		ops.ReadWriteFilesFromURL(string(body), rMethod, urlPath, w)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Not Found"))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		EntryPoint(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
