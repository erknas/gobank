package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) error {
	user := new(UserCreateRequest)

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		return err
	}

	if err := s.store.Register(r.Context(), user); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, nil)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) error {
	name := chi.URLParam(r, "name")

	user, err := s.store.Get(r.Context(), name)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, user)
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPFunc(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, nil)
		}
	}
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)

	return json.NewEncoder(w).Encode(v)
}
