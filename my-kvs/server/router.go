package server

import (
	"my-kvs/server/handlers"
	"my-kvs/server/logger"
	"net/http"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", logger.LogWrapper(handlers.Ping))
	mux.HandleFunc("/store/", logger.LogWrapper(handlers.Store))
	mux.HandleFunc("/list/", logger.LogWrapper(handlers.List))
	mux.HandleFunc("/list", logger.LogWrapper(handlers.ListAll))
	mux.HandleFunc("/shutdown", logger.LogWrapper(handlers.Shutdown))
	mux.HandleFunc("/login", logger.LogWrapper(handlers.Login))

	return mux
}
