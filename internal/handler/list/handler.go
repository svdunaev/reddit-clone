package list

import (
	"encoding/json"
	"net/http"
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

func (h *Handler) HandleGetList(w http.ResponseWriter, r *http.Request) {
	posts, err := h.repo.List(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error(), nil)
		return
	}
	if err = json.NewEncoder(w).Encode(posts); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error(), nil)
	}
}
