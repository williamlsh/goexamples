package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"goexamples/google"
	"goexamples/userip"
)

func handleSearch(w http.ResponseWriter, r *http.Request) {
	// ctx is the Context for this handler. Calling cancel closes the
	// ctx.Done channel, which is the cancellation signal for requests
	// started by this handler.
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	timeout, err := time.ParseDuration(r.FormValue("timeout"))
	if err == nil {
		// The request has a timeout, so create a context that is
		// canceled automatically when the timeout expires.
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		// The request doesn't include a timeout parameter, then cancel without timeout.
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel() // Cancel ctx as soon as handleSearch returns.

	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "no query", http.StatusBadRequest)
		return
	}
	userIP, err := userip.FromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx = userip.NewContext(ctx, userIP)

	// Run the google search and print the results.
	start := time.Now()
	results, err := google.Search(ctx, query)
	elapsed := time.Since(start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render search results.
	if err := resultsTemplate.Execute(w, struct {
		Results          google.Results
		Timeout, Elapsed time.Duration
	}{
		Results: results,
		Timeout: timeout,
		Elapsed: elapsed,
	}); err != nil {
		log.Println(err)
		return
	}
}

var resultsTemplate = template.Must(template.New("results").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Google search results</title>
  </head>
  <body>
    <ol>
      {{range .Results}}
      <li>
        {{.Title}} - <a href="{{.URL}}">{{.URL}}</a>
      </li>
      {{end}}
    </ol>
    <p>{{ len .Results }} results in {{.Elapsed}}; timeout {{.Timeout}}</p>
  </body>
</html>
`))

func main() {
	http.HandleFunc("/search", handleSearch)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
