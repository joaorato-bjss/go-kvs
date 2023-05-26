package logger

import (
	"log"
	"net/http"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func SetLogs() *os.File {
	file, err := os.OpenFile("htaccess.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return file
}

func LogWrapper(endpoint func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//log here
		InfoLogger.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		endpoint(w, r)
	}
}
