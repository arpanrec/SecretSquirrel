package iss

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
)

func TfstateHandeler(b string, m string, p string, q map[string][]string,
	w http.ResponseWriter) {

	stateFilePath := path.Join(issDataDir, p)
	lockFilePath := path.Join(issDataDir, fmt.Sprintf("%s.lock", p))

	switch m {

	case http.MethodGet:
		d, err := GetData(stateFilePath)
		if err != nil {
			if strings.HasSuffix(err.Error(), "no such file or directory") {
				HttpResponseWriter(w, http.StatusOK, "")
				return
			}
			HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", err.Error()))
			return
		}
		HttpResponseWriter(w, http.StatusOK, d)
		return
	case "LOCK":
		existingLockData, existingLockDataErr := GetData(lockFilePath)
		if existingLockDataErr != nil {
			if !strings.HasSuffix(existingLockDataErr.Error(), "no such file or directory") {
				HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", existingLockDataErr.Error()))
				return
			} else {
				_, lockDataWriteErr := PutData(lockFilePath, b)
				if lockDataWriteErr != nil {
					HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", lockDataWriteErr.Error()))
					return
				}
			}
		}

		if existingLockData != "" {
			log.Printf("Lock already exists: %s", existingLockData)
			HttpResponseWriter(w, http.StatusLocked, existingLockData)
			return
		}
	case "UNLOCK":
		err := DeleteData(lockFilePath)
		if err != nil {
			HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", err.Error()))
			return
		}
		HttpResponseWriter(w, http.StatusOK, "")
		return
	case http.MethodPut, http.MethodPost:
		if q["force"] != nil {
			if q["force"][0] == "true" {
				_, err := PutData(stateFilePath, b)
				if err != nil {
					HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", err.Error()))
					return
				}
				HttpResponseWriter(w, http.StatusOK, b)
				return
			}
		}
		_, err := PutData(stateFilePath, b)
		if err != nil {
			HttpResponseWriter(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", err.Error()))
			return
		}
		HttpResponseWriter(w, http.StatusOK, b)
		return
	default:
		HttpResponseWriter(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

}
