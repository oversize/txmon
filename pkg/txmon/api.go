package txmon

import (
	"log"
	"net/http"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	server := http.Server{
		Addr: s.addr,
		Handler: router,
	}
	log.Print("Server started")

	return server.ListenAndServe()
}
