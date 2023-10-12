package ops

import (
	"gitlab.com/arpanrecme/initsecureserver/internal/storage"
	"log"
	"net/http"
)

func ReadWriteFilesFromURL(b string, m string, filePath string, w http.ResponseWriter) {

	switch m {

	case http.MethodGet:
		d, err := storage.GetData(filePath)
		if err != nil {
			httpResponseWriter(w, http.StatusNotFound, "Not Found")
			return
		}
		httpResponseWriter(w, http.StatusOK, d)
		return
	case http.MethodPut, http.MethodPost:
		log.Println("Body: ", b)
		_, err := storage.PutData(filePath, b)
		if err != nil {
			httpResponseWriter(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		httpResponseWriter(w, http.StatusCreated, "OK")
		return
	case http.MethodDelete:
		err := storage.DeleteData(filePath)
		if err != nil {
			httpResponseWriter(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	default:
		httpResponseWriter(w, http.StatusMethodNotAllowed, "Unsupported Method")

	}
}
