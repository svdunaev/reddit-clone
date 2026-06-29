package update_post

import (
	"encoding/json"
	"errors"
	"net/http"
	"reddit-clone/internal/application/command/update_post"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/transport/http/respond"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	cmd UpdatePostCommandHandler
}

type Request struct {
	Author        string `json:"author"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	SubredditName string `json:"subreddit_name"`
}

func NewHandler(cmd UpdatePostCommandHandler) *Handler {
	return &Handler{
		cmd: cmd,
	}
}

func (h *Handler) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, IncorrectIdErrorCode, "invalid id: "+err.Error())
		return
	}

	var payload Request
	if err := decoder.Decode(&payload); err != nil {
		respond.Error(w, http.StatusBadRequest, "bad_request", "invalid JSON "+err.Error())
		return
	}

	cmd := update_post.Command{
		Id:     id,
		Author: payload.Author,
		Title:  payload.Title,
		Body:   payload.Body,
	}

	result, err := h.cmd.Handle(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respond.Error(w, http.StatusNotFound, NotFoundErrorCode, err.Error())
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			respond.ValidationFailed(w, result.Errors.Code, result.Errors.Details)
			return
		}
		respond.Error(w, http.StatusInternalServerError, "internal_error", "something went wrong")
		return
	}

	respond.JSON(w, http.StatusOK, respond.FromPost(result.Post))
}
