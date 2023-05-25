package main

import (
	"fmt"
	"my-kvs/server"
	"my-kvs/server/logger"
	"my-kvs/store"
	"net/http"
	"os"
	"strconv"
)

func main() {
	logFile := logger.SetLogs()
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {

		}
	}(logFile)

	args := os.Args

	depth := setDepth(args)
	fmt.Printf("Depth of LRU: %v\n", depth)

	port := setPort(args)

	logger.InfoLogger.Println("Spinning up server")

	fmt.Printf("Listening on %v\n", port)
	fmt.Printf("Logging to %v\n", logFile.Name())

	store.InitStore(depth)

	mux := server.RegisterRoutes()

	serv := &http.Server{Addr: port, Handler: mux}

	logger.ErrorLogger.Fatal(serv.ListenAndServe())
}

func setPort(args []string) string {
	if len(args) < 3 || args[1] != "--port" {
		logger.ErrorLogger.Fatal("exit code -1: format should be './store --port <port>'")
	}
	port, err := strconv.Atoi(args[2])
	if err != nil {
		logger.ErrorLogger.Fatalf("exit code -1: failure to parse %s into an integer port", args[2])
	}
	return fmt.Sprintf(":%d", port)
}

func setDepth(args []string) int {
	if len(args) == 5 && args[3] == "--depth" {
		depth, err := strconv.Atoi(args[4])
		if err != nil {
			logger.ErrorLogger.Fatalf("exit code -1: failure to parse %s into an integer depth", args[4])
		}
		return depth
	}

	return 0
}
