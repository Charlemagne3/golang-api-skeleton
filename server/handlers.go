package main

import (
	"net/http"
)

func (s *AppServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	if r.URL.Path != "/" {
		response = map[string]interface{}{
			"host": r.Host,
			"path": r.URL.Path,
		}
		s.writeJSON(w, r, response, http.StatusNotFound)
		return
	}
	response = map[string]interface{}{
		"host": r.Host,
		"path": r.URL.Path,
	}
	s.writeJSON(w, r, response, http.StatusOK)
}
