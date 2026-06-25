package get_post

import (
	"errors"
	"net/http"
	"reddit-clone/internal/application/query/get_post"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/transport/http/respond"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	query GetPostQueryHandler
}

func NewHandler(query GetPostQueryHandler) *Handler {
	return &Handler{
		query: query,
	}
}

func (h *Handler) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, IncorrectIdErrorCode, "invalid id: "+err.Error())
		return
	}

	query := get_post.Query{
		Id: id,
	}

	result, err := h.query.Handle(r.Context(), query)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respond.Error(w, http.StatusNotFound, result.Error.Code, err.Error())
			return
		}
		respond.Error(w, http.StatusInternalServerError, "internal_error", "something went wrong")
		return
	}

	respond.JSON(w, http.StatusOK, respond.FromPost(result.Post))
}
