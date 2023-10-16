package common

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func DeleteFileSureOrStop(l string) {
	log.Println("Deleting file: ", l)
	_, err := os.Stat(l)
	if os.IsNotExist(err) {
		log.Println("File does not exist: ", l)
	} else {
		log.Println("Deleting file: ", l)
		err := os.Remove(l)
		if err != nil {
			log.Fatalln("Error deleting file: ", err)
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
