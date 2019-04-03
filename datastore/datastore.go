package datastore

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq" // postgres
)

type Book struct {
	Isbn   string
	Title  string
	Author string
	Price  float32
}

type Datastore interface {
	AllBooks() ([]*Book, error)
}

// DB implements Datastore interface.
type DB struct {
	*sql.DB
}

// NewDB returns a new *DB which implements Datastore interface.
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) AllBooks() ([]*Book, error) {
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bks := make([]*Book, 0)
	for rows.Next() {
		bk := new(Book)
		err := rows.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
		if err != nil {
			return nil, err
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bks, nil
}
