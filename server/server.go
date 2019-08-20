package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Asuforce/odango/file"
)

type Server struct {
	endpoint string
	port     int
	hosts    []string
}

func (s *Server) Run() error {
	http.HandleFunc(s.endpoint, s.deployHandler)

	fmt.Printf("Running server on port: %d endpoint: /%s\nType Ctr-c to shutdown server.\n", s.port, s.endpoint)
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), nil); err != nil {
		return err
	}
	return nil
}

func (s *Server) deployHandler(w http.ResponseWriter, r *http.Request, file file.Base) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed.\n")
		return
	}

	commitID := r.URL.Path[len("/deploy/"):]

	file.Download(commitID)

	for _, v := range s.hosts {
		file.upload(commitID, v)
		file.unarchive(commitID, v)
	}
	fmt.Fprint(w, "Deploy success.\n")
}
