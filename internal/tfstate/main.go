package tfstate

import (
	"fmt"
	"github.com/arpanrec/secureserver/internal/common"
	"log"
	"net/http"
	"strings"

	"github.com/arpanrec/secureserver/internal/storage"
)

func TerraformStateHandler(b string, m string, p string, q map[string][]string,
	w http.ResponseWriter) {

	stateFilePath := p
	lockFilePath := fmt.Sprintf("%s.lock", p)

	switch m {

	case http.MethodGet:
		d, err := storage.GetData(stateFilePath)
		if err != nil {
			log.Println("Error while getting data: ", err)
			if strings.HasSuffix(err.Error(), "no such file or directory") {
				common.HttpResponseWriter(w, http.StatusOK, "")
				return
			}
			common.HttpResponseWriter(w, http.StatusInternalServerError,
				fmt.Sprintf("Internal Server Error: %s", err.Error()))
			return
		}
		common.HttpResponseWriter(w, http.StatusOK, d)
		return
	case "LOCK":
		existingLockData, existingLockDataErr := storage.GetData(lockFilePath)
		if existingLockDataErr != nil {
			if !strings.HasSuffix(existingLockDataErr.Error(), "no such file or directory") {
				common.HttpResponseWriter(w, http.StatusInternalServerError,
					fmt.Sprintf("Internal Server Error: %s", existingLockDataErr.Error()))
				return
			} else {
				_, lockDataWriteErr := storage.PutData(lockFilePath, b)
				if lockDataWriteErr != nil {
					common.HttpResponseWriter(w, http.StatusInternalServerError,
						fmt.Sprintf("Internal Server Error: %s", lockDataWriteErr.Error()))
					return
				}
			}
		}

		if existingLockData != "" {
			log.Printf("Lock already exists: %s", existingLockData)
			common.HttpResponseWriter(w, http.StatusLocked, existingLockData)
			return
		}
	case "UNLOCK":
		err := storage.DeleteData(lockFilePath)
		if err != nil {
			common.HttpResponseWriter(w, http.StatusInternalServerError,
				fmt.Sprintf("Internal Server Error: %s", err.Error()))
			return
		}
		common.HttpResponseWriter(w, http.StatusOK, "")
		return
	case http.MethodPut, http.MethodPost:
		if q["force"] != nil {
			if q["force"][0] == "true" {
				_, err := storage.PutData(stateFilePath, b)
				if err != nil {
					common.HttpResponseWriter(w, http.StatusInternalServerError,
						fmt.Sprintf("Internal Server Error: %s", err.Error()))
					return
				}
				common.HttpResponseWriter(w, http.StatusOK, b)
				return
			}
		}
		_, err := storage.PutData(stateFilePath, b)
		if err != nil {
			common.HttpResponseWriter(w, http.StatusInternalServerError,
				fmt.Sprintf("Internal Server Error: %s", err.Error()))
			return
		}
		common.HttpResponseWriter(w, http.StatusOK, b)
		return
	default:
		common.HttpResponseWriter(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
}
