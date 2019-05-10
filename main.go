package main

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
	create table person (
		first_name text,
		last_name text,
		email text
	);

	create table place (
		country text,
		city text null,
		telcode integer
	);
`

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}

func main() {
	db := sqlx.MustConnect("postgres", "")

	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("insert into person (first_name, last_name, email) values ($1, $2, $3)", "Jason", "Washington", "fn3nu@gmail.com")
	tx.MustExec("insert into person (first_name, last_name, email) values ($1, $2, $3)", "John", "Doe", "fn3nufs@gmail.com")
	tx.MustExec("insert into place (country, city, telcode) values ($1, $2, $3)", "United States", "New York", "1")
	tx.MustExec("insert into place (country, telcode) values ($1, $2)", "Hong Kong", "852")
	tx.MustExec("insert into place (country, telcode) values ($1, $2)", "Singapore", "835")
	tx.NamedExec("insert into person (first_name, last_name, email) values (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "fnerfni@gmail.com"})
	tx.Commit()

	// Get a slice of results.
	people := []Person{}
	db.Select(&people, "select * from person order by first_name asc")
	jason, john := people[0], people[1]
	fmt.Printf("%#v\n%#v", jason, john)

	// Get a single result.
	jason = Person{}
	err := db.Get(&jason, "select * from person where first_name=$1", "jason")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", jason)

	// if you have null fields and use SELECT *, you must use sql.Null* in your struct.
	places := []Place{}
	err = db.Select(&places, "select * from place order by telcode asc")
	if err != nil {
		panic(err)
	}
	usa, singsing, honkers := places[0], places[1], places[2]

	fmt.Printf("%#v\n%#v\n%#v\n", usa, singsing, honkers)
	// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
	// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
	// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}

	// Loop through rows using only one struct.
	place := Place{}
	rows, err := db.Queryx("select * from place")
	for rows.Next() {
		err := rows.StructScan(&place)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", place)
	}

	// Named queries, using `:name` as the bindvar.  Automatic bindvar support
	// which takes into account the dbtype based on the driverName on sqlx.Open/Connect.
	_, err = db.NamedExec(`insert into person (first_name, last_name, email) values (:first, :last, :email)`, map[string]interface{}{
		"first": "Bin",
		"last":  "Smuth",
		"email": "bensmith@allblacks.nz",
	})

	// Selects Mr. Smith from the database.
	_, err = db.NamedQuery("select * from person where first_name=:fn", map[string]interface{}{
		"fn": "Bin",
	})

	// Named queries can also use structs.  Their bind names follow the same rules
	// as the name -> db mapping, so struct fields are lowercased and the `db` tag
	// is taken into consideration.
	_, err = db.NamedQuery("select * from person where first_name=:first_name", jason)
}
