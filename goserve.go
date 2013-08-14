package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

const (
	usersRoute = "/users/"
)

func serveGetUsers(db *sql.DB) http.HandlerFunc {
	return http.NotFound
}

func serveGetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Path[len(usersRoute):])
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
	return http.NotFound
}

func servePutUser(db *sql.DB) http.HandlerFunc {
	return http.NotFound
}

func servePatchUser(db *sql.DB) http.HandlerFunc {
	return http.NotFound
}

func serveDeleteUser(db *sql.DB) http.HandlerFunc {
	return http.NotFound
}

func serveUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			serveGetUser(db)(w, r)
		case "POST":
			servePostUser(db)(w, r)
		case "PUT":
			servePutUser(db)(w, r)
		case "PATCH":
			servePatchUser(db)(w, r)
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

	http.HandleFunc(usersRoute, serveUsers(db))
	http.ListenAndServe(":4646", nil)
}
