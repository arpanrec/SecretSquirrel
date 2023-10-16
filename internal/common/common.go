package common

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func DeleteFileSureOrStop(l string) {
	_, err := os.Stat(l)
	if err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			err := os.Remove(l)
			if err != nil {
				log.Println("Error deleting file: ", err)
			}
		}
	}
}

func HttpResponseWriter(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, body)
	if err != nil {
		log.Println("Error writing response: ", err)
	}
}
