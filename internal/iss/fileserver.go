package iss

import (
	"fmt"
	"log"
	"net/http"
	"path"
)

func HttpResponseWriter(w http.ResponseWriter, code int, body string) {
	log.Printf("Response: %d %s", code, body)
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, body)
	if err != nil {
		log.Fatal(err)
	}
}
func ReadWriteFilesFromURL(b string, m string, p string, w http.ResponseWriter) {

	filePath := path.Join(issDataDir, p)

	switch m {

	case http.MethodGet:
		d, err := GetData(filePath)
		if err != nil {
			HttpResponseWriter(w, http.StatusNotFound, "Not Found")
		}
		HttpResponseWriter(w, http.StatusOK, d)
	case http.MethodPut, http.MethodPost:
		log.Println("Body: ", b)
		_, err := PutData(filePath, b)
		if err != nil {
			HttpResponseWriter(w, http.StatusInternalServerError, "Internal Server Error")
		}
		HttpResponseWriter(w, http.StatusCreated, "OK")
	}
}
