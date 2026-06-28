package list_posts

import (
	"net/http"
	"reddit-clone/internal/transport/http/respond"
)

type Handler struct {
	query ListPostsQueryHandler
}

func NewHandler(query ListPostsQueryHandler) *Handler {
	return &Handler{
		query: query,
	}
}

func (h *Handler) HandleListPosts(w http.ResponseWriter, r *http.Request) {
	result, err := h.query.Handle(r.Context())
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, result.Error.Code, "something went wrong")
		return
	}

	mappedPosts := make([]respond.Post, 0, len(result.Posts))

	for _, p := range result.Posts {
		mappedPosts = append(mappedPosts, respond.FromPost(&p))
	}

	respond.JSON(w, http.StatusOK, mappedPosts)
}
