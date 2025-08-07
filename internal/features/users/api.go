package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	requestTimeout = 10 * time.Second
	defaultLimit   = 10
	maxLimit       = 100
)

// Handler handles HTTP requests for user operations.
type Handler struct {
	service Service
}

// NewHandler creates a new user HTTP handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers user routes with the given HTTP mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.handleUsers)
	mux.HandleFunc("/users/", h.handleUserByID)
}

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) handleUserByID(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) listUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	users, err := h.service.ListUsers(ctx, limit, offset)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, map[string]any{
		"users": users,
		"count": len(users),
	})
}

func (h *Handler) createUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(ctx, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, user)
}

func (h *Handler) getUser(ctx context.Context, w http.ResponseWriter, id int64) {
	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, user)
}

func (h *Handler) updateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, id int64) {
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.service.UpdateUser(ctx, id, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, user)
}

func (h *Handler) deleteUser(ctx context.Context, w http.ResponseWriter, id int64) {
	err := h.service.DeleteUser(ctx, id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper methods for response handling.
func (h *Handler) writeJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if encodeErr := json.NewEncoder(w).Encode(data); encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrUserExists):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, ErrEmptyName), errors.Is(err, ErrInvalidEmail):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
