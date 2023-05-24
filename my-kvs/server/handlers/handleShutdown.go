package handlers

import (
	"my-kvs/server/logger"
	"net/http"
	"os"
	"time"
)

func Shutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		user := r.Header.Get("Authorization")

		if user == "" {
			logger.ErrorLogger.Printf("bad request, should provide a user under the Authorization header")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("user not provided"))
			if err != nil {
				logger.ErrorLogger.Printf("error writing 'user not provided' response: %s", err.Error())
			}
			return
		} else if user != "admin" {
			logger.ErrorLogger.Printf("forbidden operation, user is not admin")
			w.WriteHeader(http.StatusForbidden)
			_, err := w.Write([]byte("user is not admin"))
			if err != nil {
				logger.ErrorLogger.Printf("error writing 'user is not admin' response: %s", err.Error())
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			logger.ErrorLogger.Printf("error writing 'OK' response: %s", err.Error())
			return
		}
		logger.InfoLogger.Println("Connection Closed")

		go func() {
			time.Sleep(time.Millisecond)
			os.Exit(0)
		}()

	} else {
		logger.ErrorLogger.Printf("Wrong HTTP method (should be GET), not %s\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
