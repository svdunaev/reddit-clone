package create_post

import (
	"encoding/json"
	"errors"
	"net/http"
	"reddit-clone/internal/application/command/create_post"
	domain "reddit-clone/internal/domain/post"
	"reddit-clone/internal/transport/http/respond"
)

type Request struct {
	Author        string `json:"author"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	SubredditName string `json:"subreddit_name"`
}

type Handler struct {
	cmd CreatePostCommandHandler
}

func NewHandler(cmd CreatePostCommandHandler) *Handler {
	return &Handler{
		cmd: cmd,
	}
}

func (h *Handler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var payload Request

	if err := decoder.Decode(&payload); err != nil {
		respond.Error(w, http.StatusBadRequest, "bad_request", "invalid JSON "+err.Error())
		return
	}

	cmd := create_post.Command{
		Author:        payload.Author,
		Title:         payload.Title,
		Body:          payload.Body,
		SubredditName: payload.SubredditName,
	}

	result, err := h.cmd.Handle(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			response := respond.ToResponse(result)
			respond.JSON(w, http.StatusBadRequest, response)
			return
		}
		respond.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	response := respond.ToResponse(result)

	respond.JSON(w, http.StatusCreated, response)
}
