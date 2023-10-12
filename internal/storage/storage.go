package storage

import (
	"gitlab.com/arpanrecme/initsecureserver/internal/physical"
)

func getStorage(location string, data string) physical.Storage {
	return physical.FileStorage{
		Location: location,
		Data:     data,
	}
}

func GetData(location string) (string, error) {
	s := getStorage(location, "")
	return s.GetData()
}

func PutData(location string, data string) (bool, error) {
	s := getStorage(location, data)
	return s.PutData()
}

func DeleteData(location string) error {
	s := getStorage(location, "")
	return s.DeleteData()
}
