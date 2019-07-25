package main

func run() {
	logger := log.New(os.Stdout, "", 0)
	hs := setup(logger)
	logger.Printf("Running server on port: %d endpoint: /%s\nType Ctr-c to shutdown server.\n", port, endpoint)
	hs.ListenAndServe()
}

func setup(logger *log.Logger) *http.Server {
	return &http.Server{
		Addr: getAddr(),
		Handler: newServer(logWith(logger)),
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdelTimeout: 60 * time.Second,
	}
}

func getAddr() string {
	if port := config.Server.Port; port != "" {
		return ":" + port
	}

	return ":8080"
}

type Option func(*Server)

func logWith(logger *log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

type Server struct {
	mux *http.ServeMux
	logger *log.Logger
}

func newServer(options ...Option) *Server {
	s := &Server{logger: log.New(ioutil.Discard, "", 0)}

	for _, v := range options {
		v(s)
	}

	s.mux = http.NewServeMux()
	s.mux.HandleFunc(getEndpoint(), s.index)

	return s
}

func getEndpoint() string {
	if endpoint := config.Server.Endpoint; endpoint != "" {
		return "/" + endpoint + "/"
	}

	return "/deploy/"
}

func (s *Server) ServeHTTP()(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed.\n")
		return
	}

	s.log("%s, %s", r.Method, r.URL.Path)

	s.mux.ServeHTTP(w, r)
}

func (s *Server) log(format string, v ...interface{}) {
	s.logger.Printf(format+"\n", v...)
}

func (s *Server) deploy() {
	commitID := r.URL.Path[len("/deploy/"):]

	download(commitID)

	hosts := config.SSH.Hosts
	for i := range hosts {
		upload(commitID, hosts[i])
		unarchive(commitID, hosts[i])
	}
	s.log("Deploy success. filename: %s", commitID)
}
