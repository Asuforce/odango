package main

import (
	"fmt"
	"log"
	"net/http"
)

var config Config

const workDir = "/tmp/odango/" // TODO: Check when lunch

func deployHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed.\n")
		return
	}

	commitID := r.URL.Path[len("/deploy/"):]

	download(commitID)

	hosts := config.SSH.Hosts
	for i := range hosts {
		upload(commitID, hosts[i])
		unarchive(commitID, hosts[i])
	}
	fmt.Fprint(w, "Deploy success.\n")
}

func main() {
	config = Config{}

	endpoint := config.Server.Endpoint
	http.HandleFunc(endpoint, deployHandler)

	fmt.Printf("Running server on port: %d endpoint: /%s\nType Ctr-c to shutdown server.\n", config.Server.Port, endpoint)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
