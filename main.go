package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

type scriptJSON struct {
	Path string `json:"path"`
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed.\n")
	}

	s := r.URL.Path[len("/deploy/"):]

	fmt.Fprintf(w, "Env is %s !\n", s)

	body := r.Body
	defer body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, body)

	var script scriptJSON
	json.Unmarshal(buf.Bytes(), &script)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "POST script path: %s\n", script.Path)

	cmd := exec.Command(script.Path, s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command: %v, error: %v", cmd, err)
	} else {
		fmt.Printf("Output %s\n", out.String())
	}
}

func main() {
	http.HandleFunc("/deploy/", deployHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
