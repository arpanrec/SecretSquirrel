package api

import (
	"encoding/json"
	"errors"
	"github.com/arpanrec/secretsquirrel/internal/physical"
	"log"
	"net/http"
)

func KeyValue(data *[]byte, operation *string, key *string) (*physical.KVData, error) {

	switch *operation {

	case http.MethodGet:
		log.Println("http.MethodGet for KeyValue called " + *key)
		d, err := physical.Get(key, nil)
		if err != nil {
			log.Println("Error while getting data: ", err)
			return nil, err
		}
		return d, nil
	case http.MethodPut, http.MethodPost:
		var kvData physical.KVData
		errUnmarshal := json.Unmarshal(*data, &kvData)
		if errUnmarshal != nil {
			log.Println("Error while unmarshalling data from request: ", errUnmarshal)
			return nil, errUnmarshal
		}
		var saveUpdateErr error
		if *operation == http.MethodPut {
			saveUpdateErr = physical.Update(key, &kvData, nil)
		} else {
			saveUpdateErr = physical.Save(key, &kvData)
		}
		if saveUpdateErr != nil {
			return nil, saveUpdateErr
		}
		return nil, nil
	case http.MethodDelete:
		log.Println("http.MethodDelete for KeyValue called " + *key)
		err := physical.Delete(key, nil)
		if err != nil {
			log.Println("Error while deleting data: ", err)
			return nil, err
		}
		return nil, nil
	default:
		return nil, errors.New("unsupported Method")
	}
}
