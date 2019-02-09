package google

import (
	"context"
	"encoding/json"
	"net/http"

	"goexamples/userip"
)

// Results is an ordered list of search results.
type Results []Result

// A Result contains the title and URL of a search result.
type Result struct {
	Title, URL string
}

// Search sends query to Google search and returns the results.
func Search(ctx context.Context, query string) (Results, error) {
	// Prepare the Google Search API request.
	req, err := http.NewRequest("GET", "https://www.googleapis.com/customsearch/v1", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("key", "AIzaSyCZoczLD4CXD2vqOUoprl-H_Dg1cRog1zw")
	q.Set("cx", "000318443264806266275:5lnpyhajzgw")
	q.Set("q", query)

	// If ctx is carrying the user IP address, forward it to the server.
	// Google APIs use the user IP to distinguish server-initiated requests
	// from end-user requests.
	if userIP, ok := userip.FromContext(ctx); ok {
		q.Set("userip", userIP.String())
	}
	req.URL.RawQuery = q.Encode()

	var results Results
	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Parse the JSON search result.
		// https://developers.google.com/custom-search/v1/cse/list#request-body
		var data struct {
			Items []struct {
				Title        string
				FormattedURL string
			}
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		for _, res := range data.Items {
			results = append(results, Result{Title: res.Title, URL: res.FormattedURL})
		}
		return nil
	})
	return results, err
}

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	c := make(chan error, 1)
	req = req.WithContext(ctx)
	go func() {
		c <- f(http.DefaultClient.Do(req))
	}()
	select {
	case <-ctx.Done():
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
