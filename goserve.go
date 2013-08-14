package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

const (
	usersPath = "/users/"
)

func serveGetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Path[len(usersPath):])
		if err != nil {
			http.Error(w, "Couldn't parse ID from users path.", http.StatusBadRequest)
			return
		}

		var name string
		err = db.QueryRow("SELECT u.name FROM users u WHERE u.id = $1", id).Scan(&name)
		switch {
		case err == sql.ErrNoRows:
			http.NotFound(w, r)
			return
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(name))
	}
}

func servePostUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		name, ok := queryValues["name"]
		if !ok {
			http.Error(w, "Name to add must be URL encoded in POST request", http.StatusBadRequest)
			return
		}

		if len(name) != 1 {
			http.Error(w, "Only one name may be added at a time.", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO users (name) VALUES ($1)", name[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func serveDeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Path[len(usersPath):])
		if err != nil {
			http.Error(w, "Couldn't parse ID from users path.", http.StatusBadRequest)
			return
		}

		res, err := db.Exec("DELETE FROM users u WHERE u.id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func serveUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			serveGetUser(db)(w, r)
		case "POST":
			servePostUser(db)(w, r)
		case "DELETE":
			serveDeleteUser(db)(w, r)
		default:
		}
	}
}

func main() {
	db, err := sql.Open("postgres", "dbname=fido-research sslmode=disable")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	http.HandleFunc(usersPath, serveUsers(db))
	http.ListenAndServe(":4646", nil)
}
