package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *Server) handleRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := new(RegisterUserRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	user, err := NewUser(req.FirstName, req.LastName, req.Email, req.PasswordHash)
	if err != nil {
		return err
	}

	if err := s.store.Register(ctx, user); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, user)
}

func (s *Server) handleGetUserByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}

	user, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, user)
}

func (s *Server) handleGetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.GetUsers(ctx)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, UsersResponse{Users: users})
}

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

func makeHTTPFunc(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, RequestID{}, uuid.New().String())
		if err := fn(ctx, w, r); err != nil {
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
