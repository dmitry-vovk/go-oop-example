// This is a sample application demonstrating OOP implementation in Go
// Author: Dmitry Vovk <dmitry.vovk@gmail.com>
package main

import (
	"log"
	"net/http"
	"os"
	"storage"
)

const (
	logPath             = "/var/log/downloads.log"
	filesDir            = "/var/shared/"
	sessionsListUrlPath = "/sessions"
	fileDownloadPrefix  = "/file/"
	listenOn            = ":8080"
)

// Set up everything
func init() {
	setupLogger()
	storage.StoragePath = filesDir
	storage.UrlPrefix = fileDownloadPrefix
}

// Run the app
func main() {
	http.HandleFunc(fileDownloadPrefix, storage.FileDownload)
	http.HandleFunc(sessionsListUrlPath, storage.ListSessions)
	http.ListenAndServe(listenOn, nil)
}

// Try to open/create log file if set
func setupLogger() {
	if logPath != "" {
		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, (os.FileMode)(0644))
		if err == nil {
			log.SetOutput(logFile)
		} else {
			log.Printf("Could not open log file %s: %s", logPath, err)
		}
	}
}
