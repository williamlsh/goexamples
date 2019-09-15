// Reference: https://github.com/krishbhanushali/go-rest-unit-testing/blob/
// master/api.go
package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // postgres driver
)

var (
	host = ""
	port = ""
	dsn  = ""
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/api").Subrouter()
	http.ListenAndServe(port, apiRouter)
}

// GetEntries : Get All Entries
// URL : /entries
// Parameters: none
// Method: GET
// Output: JSON Encoded Entries object if found else JSON Encoded Exception.
func GetEntries(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", dsn)
	defer db.Close()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to database")
		return
	}

	var entries []entry
	rows, err := db.Query("select * from address_book;")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			eachEntry    entry
			id           int
			firstName    sql.NullString
			lastName     sql.NullString
			emailAddress sql.NullString
			phoneNumber  sql.NullString
		)

		err = rows.Scan(&id, &firstName, &lastName, &emailAddress, &phoneNumber)
		if err != nil {
			break
		}
		eachEntry.ID = id
		eachEntry.FirstName = firstName.String
		eachEntry.LastName = lastName.String
		eachEntry.EmailAddress = emailAddress.String
		eachEntry.PhoneNumber = phoneNumber.String
		entries = append(entries, eachEntry)
	}
	respondWithJSON(w, http.StatusOK, entries)
}

// GetEntryByID - Get Entry By ID
// URL : /entries?id=1
// Parameters: int id
// Method: GET
// Output: JSON Encoded Address Book Entry object if found else JSON Encoded
// Exception.
func GetEntryByID(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", dsn)
	defer db.Close()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to database")
		return
	}

	id := r.URL.Query().Get("id")
	var (
		firstName    sql.NullString
		lastName     sql.NullString
		emailAddress sql.NullString
		phoneNumber  sql.NullString
	)
	err = db.QueryRow("SELECT first_name, last_name, email_address, phone_number from address_book where id=?", id).Scan(&firstName, &lastName, &emailAddress, &phoneNumber)
	switch {
	case err == sql.ErrNoRows:
		respondWithError(w, http.StatusBadRequest, "No entry found with the id="+id)
		return
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	default:
		var eachEntry entry
		eachEntry.ID, _ = strconv.Atoi(id)
		eachEntry.FirstName = firstName.String
		eachEntry.LastName = lastName.String
		eachEntry.EmailAddress = emailAddress.String
		eachEntry.PhoneNumber = phoneNumber.String
		respondWithJSON(w, http.StatusOK, eachEntry)
	}
}

// CreateEntry - Create Entry
// URL : /entry
// Method: POST
// Body:
/*
 * {
 *	"first_name":"John",
 *	"last_name":"Doe",
 *	"email_address":"john.doe@gmail.com",
 *	"phone_number":"1234567890",
 }
*/
// Output: JSON Encoded Address Book Entry object if created else JSON Encoded
// Exception.
func CreateEntry(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", dsn)
	defer db.Close()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to database")
		return
	}

	var entry entry
	err = json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	}
	stmt, err := db.Prepare("insert into address_book (first_name, last_name, email_address, phone_number) values($1,$2,$3,$4)")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(entry.FirstName, entry.LastName, entry.EmailAddress, entry.PhoneNumber)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was problem entering the entry.")
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		id, _ := res.LastInsertId()
		entry.ID = int(id)
		respondWithJSON(w, http.StatusOK, entry)
	}
}

// UpdateEntry - Update Entry
// URL : /entry
// Method: PUT
// Body:
/*
 * {
 *	"id":1,
 *	"first_name":"Krish",
 *	"last_name":"Bhanushali",
 *	"email_address":"krishsb2405@gmail.com",
 *	"phone_number":"7798775575",
 }
*/
// Output: JSON Encoded Address Book Entry object if updated else JSON Encoded
// Exception.
func UpdateEntry(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", dsn)
	defer db.Close()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to database")
		return
	}

	var entry entry
	err = json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	}

	stmt, err := db.Prepare("update address_book set first_name=$1, last_name=$2, email_address=$3, phone_number=$4 where id=$5")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(entry.FirstName, entry.LastName, entry.EmailAddress, entry.PhoneNumber, entry.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was problem entering the entry.")
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		respondWithJSON(w, http.StatusOK, entry)
	}
}

// DeleteEntry -  Delete Entry By ID
// URL : /entries?id=1
// Parameters: int id
// Method: DELETE
// Output: JSON Encoded Address Book Entry object if found & deleted else JSON
// Encoded Exception.
func DeleteEntry(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", dsn)
	defer db.Close()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to database")
		return
	}

	id := r.URL.Query().Get("id")
	var (
		firstName    sql.NullString
		lastName     sql.NullString
		emailAddress sql.NullString
		phoneNumber  sql.NullString
	)
	err = db.QueryRow("SELECT first_name, last_name, email_address, phone_number from address_book where id=$1", id).Scan(&firstName, &lastName, &emailAddress, &phoneNumber)
	switch {
	case err == sql.ErrNoRows:
		respondWithError(w, http.StatusBadRequest, "No entry found with the id="+id)
		return
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	default:
		res, err := db.Exec("delete from address_book where id = $1", id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			return
		}
		if count == 1 {
			var eachEntry entry
			eachEntry.ID, _ = strconv.Atoi(id)
			eachEntry.FirstName = firstName.String
			eachEntry.LastName = lastName.String
			eachEntry.EmailAddress = emailAddress.String
			eachEntry.PhoneNumber = phoneNumber.String

			respondWithJSON(w, http.StatusOK, eachEntry)
			return
		}
	}
}

// UploadEntriesThroughCSV - Reads CSV, Parses the CSV and creates all the
// entries in the database
func UploadEntriesThroughCSV(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("csvFile")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong while opening the CSV.")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5
	csvData, err := reader.ReadAll()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong while parsing the CSV.")
		return
	}

	var entry entry
	for _, e := range csvData {
		if e[1] != "first_name" {
			entry.FirstName = e[1]
		}
		if e[2] != "last_name" {
			entry.LastName = e[2]
		}
		if e[3] != "email_address" {
			entry.EmailAddress = e[3]
		}
		if e[4] != "phone_number" {
			entry.PhoneNumber = e[4]
		}

		if entry.FirstName != "" && entry.LastName != "" && entry.EmailAddress != "" && entry.PhoneNumber != "" {
			b, err := json.Marshal(entry)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Something went wrong while parsing the CSV.")
				return
			}
			req, err := http.NewRequest(http.MethodPost, host+port+"api/entry", bytes.NewReader(b))
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Something went wrong while requesting to the Creation endpoint.")
				return
			}
			req.Header.Set("Content-Type", "application/json")
			cli := &http.Client{}
			resp, err := cli.Do(req)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Something went wrong while requesting to the Creation endpoint.")
				return
			}
			defer resp.Body.Close()
			if resp.Status == strconv.Itoa(http.StatusBadRequest) {
				respondWithError(w, http.StatusBadRequest, "Something went wrong while inserting.")
				return
			}
			if resp.Status == strconv.Itoa(http.StatusInternalServerError) {
				respondWithError(w, http.StatusInternalServerError, "Something went wrong while inserting.")
				return
			}
		}
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"success": "Upload successful"})
}

// DownloadEntriesToCSV - GetAllEntries, creates a CSV and downloads the CSV.
func DownloadEntriesToCSV(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get(host + port + "api/entries")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Somehow host could not be reached.")
		return
	}
	data, _ := ioutil.ReadAll(response.Body)
	var entries []entry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to unmarshal data.")
		return
	}
	b := new(bytes.Buffer)
	fileName := fmt.Sprintf("address-book-%d.csv", time.Now().Unix())
	writer := csv.NewWriter(b)
	heading := []string{"id", "first_name", "last_name", "email_address", "phone_number"}
	writer.Write(heading)

	for _, e := range entries {
		var record []string
		record = append(record, strconv.Itoa(e.ID))
		record = append(record, e.FirstName)
		record = append(record, e.LastName)
		record = append(record, e.EmailAddress)
		record = append(record, e.PhoneNumber)
		writer.Write(record)
	}
	writer.Flush()
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposion", "attachment;filename="+fileName)
	w.WriteHeader(http.StatusOK)
	w.Write(b.Bytes())
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
