package tfstate

import (
	"fmt"
	"gitlab.com/arpanrecme/initsecureserver/internal/iss"
	"gitlab.com/arpanrecme/initsecureserver/internal/iss/utils"
	"log"
	"net/http"
	"path"
	"strings"
)

func TerraformStateHandler(b string, m string, p string, q map[string][]string,
	w http.ResponseWriter) {

	stateFilePath := path.Join(iss.StorageDataDir, p)
	lockFilePath := path.Join(iss.StorageDataDir, fmt.Sprintf("%s.lock", p))

	switch m {

	case http.MethodGet:
		d, err := utils.GetData(stateFilePath)
		if err != nil {
			if strings.HasSuffix(err.Error(), "no such file or directory") {
				utils.HttpResponseWriter(w, http.StatusOK, "")
				return
			}
			utils.HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", err.Error()))
			return
		}
		utils.HttpResponseWriter(w, http.StatusOK, d)
		return
	case "LOCK":
		existingLockData, existingLockDataErr := utils.GetData(lockFilePath)
		if existingLockDataErr != nil {
			if !strings.HasSuffix(existingLockDataErr.Error(), "no such file or directory") {
				utils.HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", existingLockDataErr.Error()))
				return
			} else {
				_, lockDataWriteErr := utils.PutData(lockFilePath, b)
				if lockDataWriteErr != nil {
					utils.HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", lockDataWriteErr.Error()))
					return
				}
			}
		}

		if existingLockData != "" {
			log.Printf("Lock already exists: %s", existingLockData)
			utils.HttpResponseWriter(w, http.StatusLocked, existingLockData)
			return
		}
	case "UNLOCK":
		err := utils.DeleteData(lockFilePath)
		if err != nil {
			utils.HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", err.Error()))
			return
		}
		utils.HttpResponseWriter(w, http.StatusOK, "")
		return
	case http.MethodPut, http.MethodPost:
		if q["force"] != nil {
			if q["force"][0] == "true" {
				_, err := utils.PutData(stateFilePath, b)
				if err != nil {
					utils.HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", err.Error()))
					return
				}
				utils.HttpResponseWriter(w, http.StatusOK, b)
				return
			}
		}
		_, err := utils.PutData(stateFilePath, b)
		if err != nil {
			utils.HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", err.Error()))
			return
		}
		utils.HttpResponseWriter(w, http.StatusOK, b)
		return
	default:
		utils.HttpResponseWriter(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

}
