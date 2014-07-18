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
	UrlPrefix   string
	StoragePath string
)

func FileDownload(w http.ResponseWriter, r *http.Request) {
	remoteIp := getRemoteIpFromRequest(r)
	fName := getFileName(r)
	var start time.Time = time.Now()
	var fullFileName = StoragePath + fName
	// Check if file exists
	if _, err := os.Stat(fullFileName); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("File %s not found", fullFileName)
		return
	}
	// Create new CountingWriter
	cw := NewCountingWriter(w)
	cw.File = fName
	cw.Ip = remoteIp.String()
	cw.Start = start.Unix()
	log.Printf("Starting delivery to %s file \"%s\"", remoteIp, StoragePath+fName)
	// Register download start
	SessionsChannel <- cw
	// Perform download using out CountingWriter
	http.ServeFile(cw, r, StoragePath+fName)
	// Register download end
	SessionsChannel <- cw
	if cw.Count != 0 {
		log.Printf(
			"Download complete. %d bytes delivered in %f seconds to %s as %s",
			cw.Count,
			time.Since(start).Seconds(),
			remoteIp,
			fName,
		)
	}
}

func getRemoteIpFromRequest(r *http.Request) net.IP {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return net.ParseIP(host)
}

func getFileName(r *http.Request) string {
	uri, _ := url.QueryUnescape(r.RequestURI)
	return strings.TrimPrefix(uri, UrlPrefix)
}
