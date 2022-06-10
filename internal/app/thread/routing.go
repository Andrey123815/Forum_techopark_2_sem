package thread

import "github.com/fasthttp/router"

func SetThreadRouting(r *router.Router, h *Handlers) {
	r.GET("/api/thread/{slug_or_id}/details", h.GetThreadDetails)
	r.GET("/api/thread/{slug_or_id}/posts", h.GetThreadPosts)

	r.POST("/api/thread/{slug_or_id}/create", h.CreateNewPosts)
	r.POST("/api/thread/{slug_or_id}/vote", h.VoteThread)
	r.POST("/api/thread/{slug_or_id}/details", h.UpdateThreadDetails)
}
