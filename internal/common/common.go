package common

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func DeleteFile(l string) (bool, error) {
	err := os.Remove(l)
	if err != nil {
		log.Println("Error deleting file: ", err)
		return false, err
	}
	return true, nil
}

func HttpResponseWriter(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, body)
	if err != nil {
		log.Println("Error writing response: ", err)
	}
}
