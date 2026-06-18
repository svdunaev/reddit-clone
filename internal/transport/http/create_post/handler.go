package create_post

import (
	"encoding/json"
	"net/http"
	"reddit-clone/internal/application/command/create_post"
	"reddit-clone/internal/helpers/utils"
)

type Request struct {
	Author        string `json:"author"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	SubredditName string `json:"subreddit_name"`
}

type Handler struct {
	cmd CreateCommand
}

func NewHandler(cmd CreateCommand) *Handler {
	return &Handler{
		cmd: cmd,
	}
}

func (h *Handler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var payload Request

	if err := decoder.Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error(), nil)
		return
	}

	cmd := create_post.Command{
		Author:        payload.Author,
		Title:         payload.Title,
		Body:          payload.Body,
		SubredditName: payload.SubredditName,
	}

	createdPost, err := h.cmd.Handle(r.Context(), cmd)
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
