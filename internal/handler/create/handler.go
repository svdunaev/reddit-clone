package create

import (
	"encoding/json"
	"errors"
	"net/http"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/helpers/utils"
)

type Handler struct {
	repo PostRepository
}

func NewHandler(repo PostRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var post domain.Post

	if err := decoder.Decode(&post); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error(), nil)
		return
	}

	if errs, err := post.Validate(); err != nil {
		if errors.Is(err, domain.ErrValidation) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.WriteError(w, http.StatusBadRequest, "validation_error", "validation error", errs)
			return
		}
	}

	createdPost, err := h.repo.Create(r.Context(), post)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "creation_failed", err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdPost); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error(), nil)
		return
	}
}
