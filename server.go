package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	listenAddr string
	store      Storer
	quitch     chan os.Signal
}

func NewServer(listenAddr string, store Storer) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
		quitch:     make(chan os.Signal, 1),
	}
}

func (s *Server) Run(ctx context.Context) {
	router := chi.NewRouter()

	router.Post("/user", makeHTTPFunc(s.handleRegister))
	router.Post("/transaction", makeHTTPFunc(s.handleTransaction))
	router.Get("/user/{id}/transactions", makeHTTPFunc(s.handleGetTransactionsByUser))
	router.Get("/user/{id}", makeHTTPFunc(s.handleGetUserByID))
	router.Get("/users", makeHTTPFunc(s.handleGetUsers))

	srv := &http.Server{
		Addr:    s.listenAddr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	fmt.Printf("Server runing on [http://localhost%s]\n", s.listenAddr)

	signal.Notify(s.quitch, syscall.SIGINT, syscall.SIGTERM)
	<-s.quitch

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nServer shutdown")
}
