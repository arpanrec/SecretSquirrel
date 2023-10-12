package storage

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

func getPath() string {
	var StorageDataDir string = "data"
	issDataDirEnv := os.Getenv("ISS_DATA_DIR")
	if issDataDirEnv != "" {
		StorageDataDir = issDataDirEnv
	}
	return StorageDataDir
}

func GetData(location string) (string, error) {
	p := path.Join(getPath(), location)

	// Read the file and return the contents
	d, err := os.ReadFile(p)
	return string(d), err
}

func PutData(location string, data string) (bool, error) {
	p := path.Join(getPath(), location)
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

func DeleteData(location string) error {
	p := path.Join(getPath(), location)
	err := os.Remove(p)
	return err
}
