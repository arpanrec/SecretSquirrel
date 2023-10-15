package fileserver

import (
	"fmt"
	"github.com/arpanrec/secureserver/internal/common"
	"log"
	"net/http"
	"strings"

	"github.com/arpanrec/secureserver/internal/storage"
)

func ReadWriteFilesFromURL(b string, m string, filePath string, w http.ResponseWriter) {

	switch m {

	case http.MethodGet:
		d, err := storage.GetData(filePath)
		if err != nil {
			log.Println("Error while getting data: ", err)
			if strings.HasSuffix(err.Error(), "no such file or directory") {
				common.HttpResponseWriter(w, http.StatusNotFound, "")
				return
			}
			common.HttpResponseWriter(w, http.StatusInternalServerError,
				fmt.Sprintf("Internal Server Error: %s", err.Error()))
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
