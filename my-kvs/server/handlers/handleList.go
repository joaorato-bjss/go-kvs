package handlers

import (
	"encoding/json"
	"my-kvs/server/logger"
	"my-kvs/store"
	"net/http"
	"strings"
)

func List(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/list/")
	if r.Method == http.MethodGet {

		data := store.DoListGet(key)

		entry := data.Data
		err := data.Error

		if err != nil {
			logger.ErrorLogger.Printf(err.Error())
			w.WriteHeader(http.StatusNotFound)
			_, err2 := w.Write([]byte(err.Error()))
			if err2 != nil {
				logger.ErrorLogger.Printf("error writing response: %s", err2.Error())
			}
			return
		}

		jsonBytes, errM := json.Marshal(entry)
		if errM != nil {
			logger.ErrorLogger.Printf("error marshalling %v: %s", entry, errM)
			w.WriteHeader(http.StatusInternalServerError)
			_, err2 := w.Write([]byte("internal server error"))
			if err2 != nil {
				logger.ErrorLogger.Printf("error writing response: %s", err2.Error())
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonBytes)
		if err != nil {
			logger.ErrorLogger.Printf("error writing response %v: %s", jsonBytes, err.Error())
		}
	} else {
		logger.ErrorLogger.Printf("Wrong HTTP method (should be GET), not %s\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
