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

	router.Post("/user", makeHTTPFunc(s.handleRegister))
	router.Post("/charge", makeHTTPFunc(s.handleCharge))
	router.Post("/transfer", makeHTTPFunc(s.handleTransfer))
	router.Get("/user/{id}", makeHTTPFunc(s.handleGetUserByID))
	router.Get("/users", makeHTTPFunc(s.handleGetUsers))
	router.Delete("/user/{id}", makeHTTPFunc(s.handleDelete))

	log.Printf("server runing on [http://localhost%s]\n", s.listenAddr)

	return http.ListenAndServe(s.listenAddr, router)
}
