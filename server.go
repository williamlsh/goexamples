package x // import github.com/williamlsh/x

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	address = ":8080"
)

// cache is a two dimensional compound data type, the inner set contains
// all cached string values at each index of of the outer slice container.
// The client query which is a string slice hits every value in a set at corresponding index.
var cache Cache

type Cache struct {
	sync.Mutex
	list []Set
}

// Set containers distinct keys, it's values are irelevent.
type Set map[string]struct{}

// Query is a client request query.
type Query struct {
	Content []string `json:"content"`
}

// Result is a server response result.
type Result struct {
	Content []bool `json:"content"`
}

// Server starts an HTTPS server in background.
func Server() {
	if err := GenerateCerts(); err != nil {
		panic(err)
	}

	http.HandleFunc("/api", HandleRequest)
	err := http.ListenAndServeTLS(address, certFile, keyFile, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// HandleRequest handles /api requests.
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		http.Error(w, "content type must be application/json", http.StatusBadRequest)
		return
	}

	var query Query
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result Result
	innerResult := processQuery(query.Content, &cache)
	result.Content = innerResult

	fmt.Printf("process result: %+v\n", innerResult)

	if err := json.NewEncoder(w).Encode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func processQuery(content []string, cache *Cache) (result []bool) {
	// Prepare result first.
	result = make([]bool, len(content))

	// Empty query
	if len(content) == 0 {
		fmt.Println("Empty query content")
		return
	}

	// Lock process for this query request.
	cache.Lock()
	defer cache.Unlock()

	// For the first time client query, just copy content to cache.
	if cache.list == nil {
		cache.list = make([]Set, len(content))

		for i, s := range content {
			// Make a new Set for every cache index.
			set := make(Set)
			// Fill Set with value.
			set[s] = struct{}{}
			// Add Set to cache at corresponding index.
			cache.list[i] = set
		}

		return
	}

	// This is not the first time client query.
	for i, s := range content {
		// Case 1: client query length is longger than cache length.
		if i > len(cache.list)-1 {
			// Expand cache length to match content length.

			// Make a new Set.
			set := make(Set)
			// Fill set with new value copied from content at corresponding index.
			set[s] = struct{}{}
			// Append set to existing cache.
			cache.list = append(cache.list, set)

			continue
		}

		// Case 2: client query length is no longger than cache length.

		// Compare corresponding index value between content elements and cache Set.
		// Get Set at index i.
		set := cache.list[i]
		if _, ok := set[s]; ok {
			// Content element i just hits cache Set i.
			result[i] = true
		} else {
			// Add new element to set.
			set[s] = struct{}{}
		}
	}

	return
}
