package main

import (
	"fmt"
	"log"
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

	downloadObject(commitID)

	hosts := config.SSH.Hosts
	for i := range hosts {
		upload(commitID, hosts[i])
		unarchive(commitID, hosts[i])
	}
}

func main() {
	config = readConfig(config)

	http.HandleFunc("/deploy/", deployHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
