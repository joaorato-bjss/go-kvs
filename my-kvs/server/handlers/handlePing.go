package handlers

import (
	"my-kvs/server/logger"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain")
			logger.ErrorLogger.Println(err)
		}
		logger.InfoLogger.Println("Store pinged")
	} else {
		logger.ErrorLogger.Printf("Wrong HTTP method (should be GET), not %s\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "text/plain")
	}
}
