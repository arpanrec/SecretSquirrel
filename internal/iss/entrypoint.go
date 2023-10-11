package iss

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var issDataDir string = "data"

func EntryPoint(w http.ResponseWriter, r *http.Request) {
	issDataDirEnv := os.Getenv("ISS_DATA_DIR")
	if issDataDirEnv != "" {
		issDataDir = issDataDirEnv
	}

	urlPath := r.URL.Path
	log.Println("URL Path: ", urlPath)

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
	log.Println("Body: ", string(body))

	rMethod := r.Method
	log.Println("Method: ", rMethod)

	if strings.HasPrefix(urlPath, "ftstate/") {
		HttpResponseWriter(w, 404, "Not Found")
		return
	}

	ReadWriteFilesFromURL(string(body), rMethod, urlPath, w)

}
