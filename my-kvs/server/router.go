package server

import (
	"my-kvs/server/handlers"
	"net/http"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", handlers.Ping)
	mux.HandleFunc("/store/", handlers.Store)
	mux.HandleFunc("/list/", handlers.List)
	mux.HandleFunc("/list", handlers.ListAll)

	return mux
}
