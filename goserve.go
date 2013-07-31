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

func serveUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Path[len(usersRoute):])
		if err != nil {
			http.Error(w, "Couldn't parse ID from users path.", http.StatusBadRequest)
			return
		}

		var name string
		err = db.QueryRow("SELECT u.name FROM users u WHERE u.id = $1", id).Scan(&name)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Write([]byte(name))
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
