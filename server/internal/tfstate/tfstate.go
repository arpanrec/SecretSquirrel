package tfstate

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/arpanrec/secretsquirrel/internal/storage"
)

func TerraformStateHandler(b string, m string, p string, q map[string][]string) (int, string) {

	stateFilePath := p
	lockFilePath := fmt.Sprintf("%s.lock", p)

	switch m {

	case http.MethodGet:
		d, err := storage.GetData(stateFilePath)
		if err != nil {
			log.Println("Error while getting data: ", err)
			if strings.HasSuffix(err.Error(), "no such file or directory") {
				return http.StatusOK, ""
			}
			return http.StatusInternalServerError,
				fmt.Sprintf("Internal Server Error: %s", err.Error())
		}
		return http.StatusOK, d
	case "LOCK":
		existingLockData, existingLockDataErr := storage.GetData(lockFilePath)
		if existingLockDataErr != nil {
			if !strings.HasSuffix(existingLockDataErr.Error(), "no such file or directory") {
				return http.StatusInternalServerError,
					fmt.Sprintf("Internal Server Error: %s", existingLockDataErr.Error())
			} else {
				_, lockDataWriteErr := storage.PutData(lockFilePath, b)
				if lockDataWriteErr != nil {
					return http.StatusInternalServerError,
						fmt.Sprintf("Internal Server Error: %s", lockDataWriteErr.Error())
				}
				return http.StatusOK, ""
			}
		}

		if existingLockData != "" {
			log.Printf("Lock already exists: %s", existingLockData)
			return http.StatusLocked, existingLockData
		}
	case "UNLOCK":
		err := storage.DeleteData(lockFilePath)
		if err != nil {
			return http.StatusInternalServerError,
				fmt.Sprintf("Internal Server Error: %s", err.Error())
		}
		return http.StatusOK, ""
	case http.MethodPut, http.MethodPost:
		if q["force"] != nil {
			if q["force"][0] == "true" {
				_, err := storage.PutData(stateFilePath, b)
				if err != nil {
					return http.StatusInternalServerError,
						fmt.Sprintf("Internal Server Error: %s", err.Error())
				}
				return http.StatusOK, b
			}
		}
		_, err := storage.PutData(stateFilePath, b)
		if err != nil {
			return http.StatusInternalServerError,
				fmt.Sprintf("Internal Server Error: %s", err.Error())
		}
		return http.StatusOK, b
	default:
		return http.StatusMethodNotAllowed, "Method Not Allowed"
	}
	return http.StatusMethodNotAllowed, "Method Not Allowed"
}
