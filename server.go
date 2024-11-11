package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	listenAddr string
	store      Storer
}

func NewServer(listenAddr string, store Storer) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Run() error {
	router := chi.NewRouter()

	router.Post("/api/register", makeHTTPFunc(s.handleRegister))
	router.Get("/api/{name}", makeHTTPFunc(s.handleGet))

	log.Printf("server runing on [http://localhost%s]\n", s.listenAddr)

	return http.ListenAndServe(s.listenAddr, router)
}
