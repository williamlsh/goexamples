package main

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
)

const (
	headerAccept      = "Accept"
	headerContentType = "Content-Type"
)

// ExampleHandler is the handler we are using to demonstrate building an HTTP
// server with the jsonapi library.
type ExampleHandler struct{}

func (h *ExampleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(headerAccept) != jsonapi.MediaType {
		http.Error(w, "Unsupported media type", http.StatusUnsupportedMediaType)
	}

	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.createBlog
	case http.MethodPut:
		methodHandler = h.echoBlogs
	case http.MethodGet:
		if r.FormValue("id") != "" {
			methodHandler = h.showBlog
		} else {
			methodHandler = h.listBlogs
		}
	default:
		http.NotFound(w, r)
		return
	}

	methodHandler(w, r)
}

func (h *ExampleHandler) createBlog(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("blogs.create")

	blog := new(Blog)

	if err := jsonapiRuntime.UnmarshalPayload(r.Body, blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set(headerContentType, jsonapi.MediaType)

	if err := jsonapiRuntime.MarshalPayload(w, blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ExampleHandler) echoBlogs(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("blogs.list")

	blogs := fixtureBlogList()

	w.WriteHeader(http.StatusOK)
	w.Header().Set(headerContentType, jsonapi.MediaType)
	if err := jsonapiRuntime.MarshalPayload(w, blogs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ExampleHandler) showBlog(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonapiRuntime := jsonapi.NewRuntime().Instrument("blogs.show")

	blog := fixtureBlogCreate(intID)

	w.WriteHeader(http.StatusOK)
	w.Header().Set(headerContentType, jsonapi.MediaType)
	if err := jsonapiRuntime.MarshalPayload(w, blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ExampleHandler) listBlogs(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("blogs.list")

	blogs := fixtureBlogList()

	w.WriteHeader(http.StatusOK)
	w.Header().Set(headerContentType, jsonapi.MediaType)
	if err := jsonapiRuntime.MarshalPayload(w, blogs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
