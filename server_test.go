package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status: got (%v), want (%v)", status, http.StatusOK)
	}

	expected := "hello world!"
	if got := rr.Body.String(); got != expected {
		t.Errorf("unexpected body: got (%v), want (%v)", got, expected)
	}
}

func TestIndexHandlerNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/404", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("unexpected status: got (%v), want (%v)", status, http.StatusNotFound)
	}
}
