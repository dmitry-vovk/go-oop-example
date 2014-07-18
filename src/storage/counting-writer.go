package storage

import "net/http"

// Wrapper around http.ResponseWrapper and session to count delivered bytes
type CountingWriter struct {
	http.ResponseWriter `json:"-"`
	session
	Count int `json:"bytes"`
}

// Constructor
func NewCountingWriter(w http.ResponseWriter) *CountingWriter {
	return &CountingWriter{
		ResponseWriter: w,
		Count:          0,
	}
}

// Overridden method
func (w *CountingWriter) Write(data []byte) (n int, err error) {
	// Call parent
	n, err = w.ResponseWriter.Write(data)
	// Our additional code
	w.Count += n
	return
}
