package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Asuforce/odango/file"
)

// Server is odango server structure
type Server struct {
	endpoint string
	port     int
	hosts    []string
	file     file.File
}

// Run is running odango server
func (s *Server) Run() error {
	http.HandleFunc(s.endpoint, s.deployHandler)

	fmt.Printf("Running server on port: %d endpoint: /%s\nType Ctr-c to shutdown server.\n", s.port, s.endpoint)
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), nil); err != nil {
		return err
	}
	return nil
}

func (s *Server) deployHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed.\n")
		return
	}

	s.file.Name = r.URL.Path[len("/deploy/"):]

	if err := s.file.Download(); err != nil {
		fmt.Fprintf(w, "Faild to download tarball. Error: %v\n", err)
		return
	}

	for _, v := range s.hosts {
		if err := s.file.Upload(v); err != nil {
			fmt.Fprintf(w, "Faild to upload tarball to host %v. Error: %v\n", v, err)
			return
		}
		if err := s.file.Unarchive(v); err != nil {
			fmt.Fprintf(w, "Faild to uparchive tarball to host %v. Error: %v\n", v, err)
			return
		}
	}
	fmt.Fprint(w, "Deploy success.\n")
}
