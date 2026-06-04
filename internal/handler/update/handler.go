package update

import (
	"encoding/json"
	"errors"
	"net/http"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/helpers/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	repo PostRepository
}

func NewHandler(repo PostRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error(), nil)
		return
	}

	var payload domain.Post
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error(), nil)
		return
	}

	if errs, err := payload.Validate(); err != nil {
		if errors.Is(err, domain.ErrValidation) {
			utils.WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), errs)
			return
		}
	}

	post, err := h.repo.Update(r.Context(), id, payload)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, "not_found", err.Error(), nil)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error(), nil)
		return
	}

	if err := json.NewEncoder(w).Encode(post); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error(), nil)
	}
}
