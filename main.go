package main

import (
	"fmt"
	"net/http"
)

var config gongchaConfig

const workDir = "/tmp/gongcha/" // TODO: Check when lunch gongcha

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
	config = readConfig(config)
	commitID := "80d712d3fef760aa346985a837efdb37bb56cef0"
	downloadObject(commitID)
	upload(commitID)

	// http.HandleFunc("/deploy/", deployHandler)

	// log.Fatal(http.ListenAndServe(":8080", nil))
}
