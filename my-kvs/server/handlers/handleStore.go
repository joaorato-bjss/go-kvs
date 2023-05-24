package handlers

import (
	"my-kvs/server/logger"
	"net/http"
)

func Store(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		//response := store.DoStorePut()
	case http.MethodGet:

	case http.MethodDelete:

	default:
		logger.ErrorLogger.Printf("Wrong HTTP method (should be PUT, GET or DELETE), not %s\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
