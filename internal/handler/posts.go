package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/helpers"
	"reddit-clone/internal/service/posts"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	repo posts.PostRepository
}

func NewHandler(repo posts.PostRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var post domain.Post

	if err := decoder.Decode(&post); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}

	if errs, err := post.Validate(); err != nil {
		if errors.Is(err, domain.ErrValidation) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"code":    "validation_error",
					"message": "validation failed",
					"details": errs,
				},
			})
			return
		}
	}

	createdPost, err := h.repo.Create(r.Context(), post)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, "creation_failed", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdPost); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
}

func (h *Handler) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "incorrect_id", err.Error())
		return
	}

	post, err := h.repo.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			helpers.WriteError(w, http.StatusNotFound, "not_found", err.Error())
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(post); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error())
	}
}
