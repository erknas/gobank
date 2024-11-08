package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	listenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Run() error {
	router := chi.NewRouter()

	return http.ListenAndServe(s.listenAddr, router)
}
