package api

import (
	"fmt"
	"net/http"
)

type MuxFabric interface {
	InitMux() *http.ServeMux
}

type Server struct {
	S *http.Server
}

func InitServer(addr string, MxF MuxFabric) (*Server, error) {
	if MxF == nil {
		return nil, fmt.Errorf("InitServer: MuxFabric is nil")
	}

	var s Server
	s.S = &http.Server{Addr: addr}
	s.S.Handler = MxF.InitMux()
	return &s, nil
}

func (s Server) Start() error {
	err := s.S.ListenAndServe()

	return fmt.Errorf("Start: %s", err)
}
