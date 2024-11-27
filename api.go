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
		return InvalidJSON()
	}
	defer r.Body.Close()

	if errors := req.ValidateUserData(); len(errors) > 0 {
		return InvalidRequestData(errors)
	}

	user, err := NewUser(req.FirstName, req.LastName, req.Email, req.PhoneNumber, req.Password)
	if err != nil {
		return err
	}

	id, err := s.store.Register(ctx, user)
	if err != nil {
		return err
	}

	resp := RegisterUserResponse{
		StatusCode: http.StatusOK,
		Msg:        "user successfully registered",
		User: User{
			ID:          id,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			Number:      user.Number,
			Balance:     user.Balance,
			CreatedAt:   user.CreatedAt,
		},
	}

	return writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleCharge(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := new(ChargeRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return InvalidJSON()
	}
	defer r.Body.Close()

	if errors := req.ValidateChargeData(); len(errors) > 0 {
		return InvalidRequestData(errors)
	}

	balance, err := s.store.Charge(ctx, req)
	if err != nil {
		return err
	}

	resp := ChargeResponse{
		StatusCode: http.StatusCreated,
		Msg:        "successful transaction",
		Amount:     req.Amount, Balance: balance,
	}

	return writeJSON(w, http.StatusCreated, resp)
}

func (s *Server) handleTransfer(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := new(TransferRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return InvalidJSON()
	}
	defer r.Body.Close()

	if errors := req.ValidateTransferData(); len(errors) > 0 {
		return InvalidRequestData(errors)
	}

	balance, err := s.store.Transfer(ctx, req)
	if err != nil {
		return err
	}

	resp := TransferResponse{
		StatusCode: http.StatusCreated,
		Msg:        "successful transaction",
		Amount:     req.Amount, Balance: balance,
	}

	return writeJSON(w, http.StatusCreated, resp)
}

func (s *Server) handleGetUserByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return InvalidID()
	}

	user, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	resp := UserResponse{
		StatusCode: http.StatusOK,
		User:       *user,
	}

	return writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleGetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.GetUsers(ctx)
	if err != nil {
		return err
	}

	resp := UsersResponse{
		StatusCode: http.StatusOK,
		Users:      users,
	}

	return writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return InvalidID()
	}

	if err := s.store.Delete(ctx, id); err != nil {
		return err
	}

	resp := DeleteUserResponse{
		StatusCode: http.StatusOK,
		Msg:        "user successfully deleted",
		ID:         id,
	}

	return writeJSON(w, http.StatusOK, resp)
}

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

func makeHTTPFunc(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(context.Background(), RequestID{}, uuid.New().String())
		if err := fn(ctx, w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				writeJSON(w, apiErr.StatusCode, apiErr)
			} else {
				errResp := map[string]any{
					"statusCode": http.StatusInternalServerError,
					"msg":        "internal server error",
				}
				writeJSON(w, http.StatusInternalServerError, errResp)
			}
		}
	}
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}

func parseID(r *http.Request) (int, error) {
	strID := chi.URLParam(r, "id")
	return strconv.Atoi(strID)
}
