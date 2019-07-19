package main

import (
	"fmt"
	"log"
	"net/http"
)

func deployHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed.\n")
	}

	s := r.URL.Path[len("/deploy/"):]

	fmt.Fprintf(w, "Env is %s !\n", s)
}

func main() {
	http.HandleFunc("/deploy/", deployHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
