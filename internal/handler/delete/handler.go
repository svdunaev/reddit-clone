package delete

import (
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

func (h *Handler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error(), nil)
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, "not_found", err.Error(), nil)
		}
		utils.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
