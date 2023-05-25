package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"my-kvs/server/logger"
	"my-kvs/store"
	"net/http"
	"strings"
)

func Store(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/store/")
	switch r.Method {
	case http.MethodPut:
		handleStorePut(key, w, r)
	case http.MethodGet:
		handleStoreGet(key, w)
	case http.MethodDelete:
		handleStoreDelete(key, w, r)
	default:
		logger.ErrorLogger.Printf("Wrong HTTP method (should be PUT, GET or DELETE), not %s\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleStorePut(key string, w http.ResponseWriter, r *http.Request) {

	valid, user := Authorise(w, r)
	if !valid {
		return
	}

	// obtain value
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.ErrorLogger.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, err2 := w.Write([]byte(err.Error()))
		if err2 != nil {
			logger.ErrorLogger.Printf("error writing response: %s", err2.Error())
		}
		return
	}

	jsonBody, errM := json.Marshal(string(body))
	if errM != nil {
		logger.ErrorLogger.Printf("error marshalling %s: %s", string(body), errM)
		w.WriteHeader(http.StatusBadRequest)
		_, err2 := w.Write([]byte(err.Error()))
		if err2 != nil {
			logger.ErrorLogger.Printf("error writing response: %s", err2.Error())
		}
		return
	}

	var newValue any
	errUM := json.Unmarshal(jsonBody, &newValue)
	if errUM != nil {
		logger.ErrorLogger.Printf("error unmarshalling %s: %s", string(jsonBody), errUM)
		w.WriteHeader(http.StatusBadRequest)
		_, err2 := w.Write([]byte(errUM.Error()))
		if err2 != nil {
			logger.ErrorLogger.Printf("error writing response: %s", err2.Error())
		}
		return
	}

	response := store.DoStorePut(key, user, newValue)

	if response.Error != nil {
		logger.ErrorLogger.Printf(response.Error.Error())
		w.WriteHeader(http.StatusForbidden)
		_, err2 := w.Write([]byte(response.Error.Error()))
		if err2 != nil {
			logger.ErrorLogger.Printf("error writing response: %s", err2.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err2 := w.Write([]byte("value inserted with key: " + key))
	if err2 != nil {
		logger.ErrorLogger.Printf("error writing response: %s", err2.Error())
	}
}

func handleStoreGet(key string, w http.ResponseWriter) {
	data := store.DoStoreGet(key)

	if data.Error != nil {
		logger.ErrorLogger.Printf(data.Error.Error())
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(data.Error.Error()))
		if err != nil {
			logger.ErrorLogger.Printf("error writing response: %s", err.Error())
		}
	} else {
		var jsonBytes []byte
		var err error

		entryValue := data.Data
		if strEntryValue, ok := entryValue.(string); ok {
			rawData := json.RawMessage(strEntryValue)
			jsonBytes, err = rawData.MarshalJSON()
			if err != nil {
				logger.ErrorLogger.Print(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				_, err2 := w.Write([]byte("internal server error"))
				if err2 != nil {
					logger.ErrorLogger.Printf("error writing response: %s", err2.Error())
				}
				return
			}
		} else {
			jsonBytes, err = json.Marshal(entryValue)
			if err != nil {
				logger.ErrorLogger.Print(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				_, err2 := w.Write([]byte("internal server error"))
				if err2 != nil {
					logger.ErrorLogger.Printf("error writing response: %s", err2.Error())
				}
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonBytes)
		if err != nil {
			logger.ErrorLogger.Printf("error writing response with %s: %s", string(jsonBytes), err.Error())
		}
	}
}

func handleStoreDelete(key string, w http.ResponseWriter, r *http.Request) {

	valid, user := Authorise(w, r)
	if !valid {
		return
	}

	response := store.DoStoreDelete(key, user)

	err := response.Error
	if err != nil {
		if errors.Is(err, store.ErrNotOwner) {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		_, err2 := w.Write([]byte(err.Error()))
		if err2 != nil {
			logger.ErrorLogger.Printf("error writing response: %s", err2.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("entry deleted with key: " + key))
	if err != nil {
		logger.ErrorLogger.Printf("error writing response: %s", err.Error())
	}
}
