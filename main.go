package main

import (
	"context"
	"database/sql"
	"fmt"
	"goexamples/datastore"
	"log"
	"net/http"
)

type key int

const contextDbKey key = 0

type ContextInjector struct {
	ctx context.Context
	h   http.Handler
}

func (ci ContextInjector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ci.h.ServeHTTP(w, r.WithContext(ci.ctx))
}

func main() {
	db, err := datastore.NewDB("postgres://user:pass!@localhost/bookstore")
	if err != nil {
		log.Panic(err)
	}

	ctx := context.WithValue(context.Background(), contextDbKey, db)

	http.Handle("/books", ContextInjector{ctx, http.HandlerFunc(booksIndex)})
	http.ListenAndServe(":3000", nil)
}

func booksIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	db, ok := r.Context().Value(contextDbKey).(*sql.DB)
	if !ok {
		http.Error(w, "could not get database connection pool from context", 500)
		return
	}

	bks, err := datastore.AllBooks(db)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for _, bk := range bks {
		fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.Isbn, bk.Title, bk.Author, bk.Price)
	}
}
