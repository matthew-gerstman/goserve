package main

import (
	"fmt"
	"net/http"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Internet!")
}

func main() {
	http.HandleFunc("/", sayHello)
	http.ListenAndServe(":4646", nil)
}
