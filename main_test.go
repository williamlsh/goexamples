package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// statusHandler is an http.Handler that writes an empty response using itself
// as the response status code.
type statusHandler int

func (h *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(int(*h))
}

func TestIsTagged(t *testing.T) {
	// Set up a fake "Google Code" web server reporting 404 not found.
	status := statusHandler(http.StatusNotFound)
	s := httptest.NewServer(&status)
	defer s.Close()

	if isTagged(s.URL) {
		t.Fatal("isTagged == true, want false")
	}

	// Change fake server status to 200 OK and try again.
	status = http.StatusOK

	if !isTagged(s.URL) {
		t.Fatal("isTagged == false, want true")
	}
}
