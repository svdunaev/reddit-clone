package delete_post

import (
	"errors"
	"net/http"
	"reddit-clone/internal/application/command/delete_post"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/transport/http/respond"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	cmd DeletePostCommandHandler
}

func NewHandler(cmd DeletePostCommandHandler) *Handler {
	return &Handler{
		cmd: cmd,
	}
}

func (h *Handler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, IncorrectIdErrorCode, "invalid id: "+err.Error())
		return
	}

	cmd := delete_post.Command{
		Id: id,
	}

	err = h.cmd.Handle(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respond.Error(w, http.StatusNotFound, NotFoundErrorCode, err.Error())
			return
		}
		respond.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}
