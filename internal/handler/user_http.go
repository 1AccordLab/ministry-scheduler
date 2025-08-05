package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ministry-scheduler/internal/domain"
	"ministry-scheduler/internal/usecase"
)

const (
	requestTimeout = 10 * time.Second
	defaultLimit   = 10
	maxLimit       = 100
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(usecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.handleUsers)
	mux.HandleFunc("/users/", h.handleUserByID)
}

func (h *UserHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		h.listUsers(ctx, w, r)
	case http.MethodPost:
		h.createUser(ctx, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	if idStr == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUser(ctx, w, id)
	case http.MethodPut:
		h.updateUser(ctx, w, r, id)
	case http.MethodDelete:
		h.deleteUser(ctx, w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) listUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	users, err := h.usecase.ListUsers(ctx, limit, offset)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, map[string]any{
		"users": users,
		"count": len(users),
	})
}

func (h *UserHandler) createUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.usecase.CreateUser(ctx, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, user)
}

func (h *UserHandler) getUser(ctx context.Context, w http.ResponseWriter, id int64) {
	user, err := h.usecase.GetUser(ctx, id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, user)
}

func (h *UserHandler) updateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, id int64) {
	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.usecase.UpdateUser(ctx, id, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, user)
}

func (h *UserHandler) deleteUser(ctx context.Context, w http.ResponseWriter, id int64) {
	err := h.usecase.DeleteUser(ctx, id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper methods for response handling.
func (h *UserHandler) writeJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if encodeErr := json.NewEncoder(w).Encode(data); encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrUserExists):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, domain.ErrEmptyName), errors.Is(err, domain.ErrInvalidEmail):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
