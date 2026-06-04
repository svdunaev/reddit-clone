package get_by_id

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

func (h *Handler) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "incorrect_id", err.Error(), nil)
		return
	}

	post, err := h.repo.GetById(r.Context(), id)
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
