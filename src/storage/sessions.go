// Here we handle download sessions
package storage

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON serializable structure describing a session
type session struct {
	File  string `json:"file"`
	Start int64  `json:"start"`
	Ip    string `json:"ip"`
}

// JSON-able list of all current sessions
type sessionsList struct {
	Sessions []*CountingWriter `json:"sessions"`
}

var (
	sessions        = make(map[*CountingWriter]bool)
	SessionsChannel = make(chan *CountingWriter)
)

// Run session tracker
func init() {
	go track(SessionsChannel)
}

// Session tracker
func track(s chan *CountingWriter) {
	// Read from the channel until it closed
	for cw := range s {
		// Check if we already got the session
		if _, has := sessions[cw]; has {
			// If yes, then we shut it down
			delete(sessions, cw)
		} else {
			// If no, then this is new session
			sessions[cw] = true
		}
	}
}

// HTTP endpoint that lists current sessions in JSON
func ListSessions(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-type", "application/json; charset=UTF-8")
	var s sessionsList
	for session := range sessions {
		s.Sessions = append(s.Sessions, session)
	}
	// Encode sessions list into JSON and output it
	if out, err := json.MarshalIndent(s, "", "  "); err == nil {
		w.Write(out)
	} else {
		log.Printf("Could not marshal response: %s", err)
	}
}
