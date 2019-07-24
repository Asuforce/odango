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

	download(commitID)

	hosts := config.SSH.Hosts
	for i := range hosts {
		upload(commitID, hosts[i])
		unarchive(commitID, hosts[i])
	}
	fmt.Fprint(w, "Deploy success.\n")
}

func main() {
	config = readConfig(config)

	port := config.Server.Port
	if port == 0 {
		port = 8080
	}

	endpoint := config.Server.Endpoint
	if endpoint == "" {
		endpoint = "deploy"
	}
	http.HandleFunc("/"+endpoint+"/", deployHandler)

	fmt.Printf("Running server on port: %d endpoint: /%s\nType Ctr-c to shutdown server.\n", port, endpoint)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
