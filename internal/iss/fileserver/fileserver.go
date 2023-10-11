package fileserver

import (
	"gitlab.com/arpanrecme/initsecureserver/internal/iss"
	"gitlab.com/arpanrecme/initsecureserver/internal/iss/utils"
	"log"
	"net/http"
	"path"
)

func ReadWriteFilesFromURL(b string, m string, p string, w http.ResponseWriter) {

	filePath := path.Join(iss.StorageDataDir, p)

	switch m {

	case http.MethodGet:
		d, err := utils.GetData(filePath)
		if err != nil {
			utils.HttpResponseWriter(w, http.StatusNotFound, "Not Found")
			return
		}
		utils.HttpResponseWriter(w, http.StatusOK, d)
		return
	case http.MethodPut, http.MethodPost:
		log.Println("Body: ", b)
		_, err := utils.PutData(filePath, b)
		if err != nil {
			utils.HttpResponseWriter(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		utils.HttpResponseWriter(w, http.StatusCreated, "OK")
		return
	case http.MethodDelete:
		err := utils.DeleteData(filePath)
		if err != nil {
			utils.HttpResponseWriter(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	default:
		utils.HttpResponseWriter(w, http.StatusMethodNotAllowed, "Unsupported Method")

	}
}
