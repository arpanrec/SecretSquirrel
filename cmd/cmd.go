package cmd

import (
	"fmt"
	"github.com/arpanrec/secureserver/internal/auth"
	"github.com/arpanrec/secureserver/internal/common"
	"github.com/arpanrec/secureserver/internal/fileserver"
	"github.com/arpanrec/secureserver/internal/tfstate"
	"io"
	"log"
	"net/http"
	"strings"
)

func entryPoint(w http.ResponseWriter, r *http.Request) {

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

	authHeader := header.Get("Authorization")
	username, ok := auth.GetUserDetails(authHeader)
	if !ok {
		common.HttpResponseWriter(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	log.Println("URL Path: ", urlPath, "\nMethod: ", rMethod, "\nHeader: ", header,
		"\nForm Data: ", formData,
		"\nBody: ", string(body), "\nQuery: ", query)

	locationPath := fmt.Sprintf("%v/%v", username, urlPath[3:])

	if strings.HasPrefix(urlPath, "/v1/tfstate/") {
		tfstate.TerraformStateHandler(string(body), rMethod, locationPath, query, w)
	} else if strings.HasPrefix(urlPath, "/v1/files/") {
		fileserver.ReadWriteFilesFromURL(string(body), rMethod, locationPath, w)
	} else {
		common.HttpResponseWriter(w, http.StatusNotFound, "Not Found")
	}
}

func Runner() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		entryPoint(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
