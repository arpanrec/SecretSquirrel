package iss

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var StorageDataDir string = "data"

func EntryPoint(w http.ResponseWriter, r *http.Request) {
	issDataDirEnv := os.Getenv("ISS_DATA_DIR")
	if issDataDirEnv != "" {
		StorageDataDir = issDataDirEnv
	}

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
		TerraformStateHandler(string(body), rMethod, urlPath, query, w)
	} else if strings.HasPrefix(urlPath, "/files/") {
		ReadWriteFilesFromURL(string(body), rMethod, urlPath, w)
	} else {
		HttpResponseWriter(w, http.StatusNotFound, "")
	}
}
