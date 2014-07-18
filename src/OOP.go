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
)

func init() {
	setupLogger()
	storage.StoragePath = filesDir
	storage.UrlPrefix = fileDownloadPrefix
}

func main() {
	http.HandleFunc(fileDownloadPrefix, storage.FileDownload)
	http.HandleFunc(sessionsListUrlPath, storage.ListSessions)
	http.ListenAndServe(":8080", nil)
}

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
