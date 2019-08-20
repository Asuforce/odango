package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Asuforce/odango/file"
)

// Server is odango server structure
type Server struct {
	Endpoint string
	Port     int
	Hosts    []string
	File     file.File
}

// Run is running odango server
func (s *Server) Run() error {
	http.HandleFunc(s.Endpoint, s.deployHandler)

	fmt.Printf("Running server on port: %d endpoint: /%s\nType Ctr-c to shutdown server.\n", s.Port, s.Endpoint)
	if err := http.ListenAndServe(":"+strconv.Itoa(s.Port), nil); err != nil {
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

	s.File.Name = r.URL.Path[len("/deploy/"):]

	if err := s.File.Download(); err != nil {
		fmt.Fprintf(w, "Faild to download tarball. Error: %v\n", err)
		return
	}

	for _, v := range s.Hosts {
		if err := s.File.Upload(v); err != nil {
			fmt.Fprintf(w, "Faild to upload tarball to host %v. Error: %v\n", v, err)
			return
		}
		if err := s.File.Unarchive(v); err != nil {
			fmt.Fprintf(w, "Faild to uparchive tarball to host %v. Error: %v\n", v, err)
			return
		}
	}
	fmt.Fprint(w, "Deploy success.\n")
}
