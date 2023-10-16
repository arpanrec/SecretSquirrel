package fileserver

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/arpanrec/secureserver/internal/storage"
)

func ReadWriteFilesFromURL(b string, m string, filePath string) (int, string) {

	switch m {

	case http.MethodGet:
		d, err := storage.GetData(filePath)
		if err != nil {
			log.Println("Error while getting data: ", err)
			if strings.HasSuffix(err.Error(), "no such file or directory") {
				return http.StatusNotFound, ""
			}
			return http.StatusInternalServerError,
				fmt.Sprintf("Internal Server Error: %s", err.Error())
		}
		return http.StatusOK, d
	case http.MethodPut, http.MethodPost:
		log.Println("Body: ", b)
		_, err := storage.PutData(filePath, b)
		if err != nil {
			return http.StatusInternalServerError, "Internal Server Error"
		}
		return http.StatusCreated, "OK"
	case http.MethodDelete:
		err := storage.DeleteData(filePath)
		if err != nil {
			return http.StatusInternalServerError, "Internal Server Error"
		}
	default:
		return http.StatusMethodNotAllowed, "Unsupported Method"
	}
	return http.StatusMethodNotAllowed, "Unsupported Method"
}
