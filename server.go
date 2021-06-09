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

// cache is a concurency safe set.
var cache Cache

type Cache struct {
	sync.Mutex
	Set
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
	if cache.Set == nil {
		// Initialize set.
		cache.Set = make(Set)

		// Fill set with distinct query elements.
		for _, s := range content {
			cache.Set[s] = struct{}{}
		}
		return
	}

	// This is not the first time client query.

	// Copy a set from cache, use this set as base cache to compare user query,
	// so duplicated query elements won't effect final results.
	set := make(Set)
	for k, v := range cache.Set {
		set[k] = v
	}

	for i, s := range content {
		// Set contains query elememt.
		if _, ok := set[s]; ok {
			result[i] = true
		} else {
			// Query elememt is not in set.
			cache.Set[s] = struct{}{}
		}
	}

	return
}
