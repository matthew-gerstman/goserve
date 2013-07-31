package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
)

func sayHello(name string) http.HandlerFunc {
	helloPhrase := fmt.Sprintf("Hello, %v!", name)

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, helloPhrase)
	}
}

func main() {
	db, err := sql.Open("postgres", "dbname=fido-research sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer db.Close()

	var name string
	db.QueryRow("SELECT name FROM users").Scan(&name)

	http.HandleFunc("/", sayHello(name))
	http.ListenAndServe(":4646", nil)
}
