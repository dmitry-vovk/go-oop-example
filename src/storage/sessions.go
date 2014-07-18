package storage

import (
	"encoding/json"
	"log"
	"net/http"
)

type session struct {
	File  string `json:"file"`
	Start int64  `json:"start"`
	Ip    string `json:"ip"`
}

type sessionsList struct {
	Sessions []*CountingWriter `json:"sessions"`
}

var (
	sessions        = make(map[*CountingWriter]bool)
	SessionsChannel = make(chan *CountingWriter)
)

func init() {
	go track(SessionsChannel)
}

func track(s chan *CountingWriter) {
	for cw := range s {
		if _, has := sessions[cw]; has {
			delete(sessions, cw)
		} else {
			sessions[cw] = true
		}
	}
}

func ListSessions(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-type", "application/json; charset=UTF-8")
	var s sessionsList
	for session := range sessions {
		s.Sessions = append(s.Sessions, session)
	}
	if out, err := json.MarshalIndent(s, "", "  "); err == nil {
		w.Write(out)
	} else {
		log.Printf("Could not marshal response: %s", err)
	}
}
