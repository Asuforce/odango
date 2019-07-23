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
		return
	}

	commitID := r.URL.Path[len("/deploy/"):]
	fmt.Fprintf(w, "CommitID is %s !\n", commitID)

	downloadObject(commitID)
}

func main() {
	http.HandleFunc("/deploy/", deployHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
