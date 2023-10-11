package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func HttpResponseWriter(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, body)
	if err != nil {
		log.Fatal(err)
	}
}

func GetData(p string) (string, error) {

	// Read the file and return the contents
	d, err := os.ReadFile(p)
	return string(d), err
}

func PutData(p string, data string) (bool, error) {
	dir := filepath.Dir(p)
	errMakeDir := os.MkdirAll(dir, 0755)
	if errMakeDir != nil {
		log.Fatal(errMakeDir)
		return false, errMakeDir
	}
	errWriteFile := os.WriteFile(p, []byte(data), 0644)
	if errWriteFile != nil {
		log.Fatal(errWriteFile)
		return false, errWriteFile
	}
	return true, nil
}

func DeleteData(p string) error {
	err := os.Remove(p)
	return err
}
