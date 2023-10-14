package fileserver

import (
	"github.com/arpanrec/secureserver/internal/common"
	"log"
	"net/http"

	"github.com/arpanrec/secureserver/internal/storage"
)

func ReadWriteFilesFromURL(b string, m string, filePath string, w http.ResponseWriter) {

	switch m {

	case http.MethodGet:
		d, err := storage.GetData(filePath)
		if err != nil {
			common.HttpResponseWriter(w, http.StatusNotFound, "Not Found")
			return
		}
		common.HttpResponseWriter(w, http.StatusOK, d)
		return
	case http.MethodPut, http.MethodPost:
		log.Println("Body: ", b)
		_, err := storage.PutData(filePath, b)
		if err != nil {
			common.HttpResponseWriter(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		common.HttpResponseWriter(w, http.StatusCreated, "OK")
		return
	case http.MethodDelete:
		err := storage.DeleteData(filePath)
		if err != nil {
			common.HttpResponseWriter(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	default:
		common.HttpResponseWriter(w, http.StatusMethodNotAllowed, "Unsupported Method")

	}
}
