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
	req := new(NewUserRequest)

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

	resp := NewUserResponse{
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

func (s *Server) handleTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	transactionType := chi.URLParam(r, "type")

	req := new(TransactionRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return InvalidJSON()
	}
	defer r.Body.Close()

	req.Type = transactionType

	if errors := req.ValidateTransaction(); len(errors) > 0 {
		return InvalidRequestData(errors)
	}

	if transactionType == "charge" {
		transaction, err := s.store.Charge(ctx, req)
		if err != nil {
			return err
		}

		resp := TransactionResponse{
			StatusCode:  http.StatusCreated,
			Msg:         "successful transaction",
			Transaction: *transaction,
		}

		return writeJSON(w, http.StatusCreated, resp)
	}

	if transactionType == "transfer" {
		transaction, err := s.store.Transfer(ctx, req)
		if err != nil {
			return err
		}

		resp := TransactionResponse{
			StatusCode:  http.StatusCreated,
			Msg:         "successful transaction",
			Transaction: *transaction,
		}

		return writeJSON(w, http.StatusCreated, resp)
	}

	return nil
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

func (s *Server) handleGetTransactions(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	transactions, err := s.store.GetTransactions(ctx)
	if err != nil {
		return err
	}

	resp := TransactionsResponse{
		StatusCode:   http.StatusOK,
		Transactions: transactions,
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
