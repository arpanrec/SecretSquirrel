package iss

import (
	"log"
	"net/http"
	"path"
)

func ReadWriteFilesFromURL(b string, m string, p string, w http.ResponseWriter) {

	filePath := path.Join(StorageDataDir, p)

	switch m {

	case http.MethodGet:
		d, err := GetData(filePath)
		if err != nil {
			HttpResponseWriter(w, http.StatusNotFound, "Not Found")
			return
		}
		HttpResponseWriter(w, http.StatusOK, d)
		return
	case http.MethodPut, http.MethodPost:
		log.Println("Body: ", b)
		_, err := PutData(filePath, b)
		if err != nil {
			HttpResponseWriter(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		HttpResponseWriter(w, http.StatusCreated, "OK")
		return
	case http.MethodDelete:
		err := DeleteData(filePath)
		if err != nil {
			HttpResponseWriter(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	default:
		HttpResponseWriter(w, http.StatusMethodNotAllowed, "Unsupported Method")

	}
}
