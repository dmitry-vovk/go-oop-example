package storage

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	UrlPrefix   string // Part of url to remove to get clean file name
	StoragePath string // Where files are located in filesystem
)

func FileDownload(w http.ResponseWriter, r *http.Request) {
	var fileName string = getFileName(r)
	var start time.Time = time.Now()
	var fullFileName = StoragePath + fileName
	var remoteIp string = getRemoteIpFromRequest(r)
	// Check if file exists
	if _, err := os.Stat(fullFileName); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("File %s not found", fullFileName)
		return
	}
	// Create new CountingWriter
	cw := NewCountingWriter(w)
	cw.File = fileName
	cw.Ip = remoteIp.String()
	cw.Start = start.Unix()
	log.Printf("Starting delivery to %s file \"%s\"", remoteIp, fullFileName)
	// Register download start
	SessionsChannel <- cw
	// Perform download using our CountingWriter
	http.ServeFile(cw, r, fullFileName)
	// Register download end
	SessionsChannel <- cw
	log.Printf(
		"Download complete. %d bytes delivered in %f seconds to %s as %s",
		cw.Count,
		time.Since(start).Seconds(),
		remoteIp,
		fileName,
	)
}

func getRemoteIpFromRequest(r *http.Request) net.IP {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return net.ParseIP(host)
}

func getFileName(r *http.Request) string {
	uri, _ := url.QueryUnescape(r.RequestURI)
	return strings.TrimPrefix(uri, UrlPrefix)
}
