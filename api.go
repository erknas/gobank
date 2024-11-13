package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) error {
	req := new(RegisterUserRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	user, err := NewUser(req.FirstName, req.LastName, req.Email, req.PasswordHash)
	if err != nil {
		return err
	}

	if err := s.store.Register(r.Context(), user); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, user)
}

func (s *Server) handleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}

	user, err := s.store.GetUserByID(r.Context(), id)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, user)
}

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.GetUsers(r.Context())
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, users)
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

func parseID(r *http.Request) (int, error) {
	id := chi.URLParam(r, "id")

	return strconv.Atoi(id)
}
