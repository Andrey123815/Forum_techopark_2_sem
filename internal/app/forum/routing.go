package forum

import "github.com/fasthttp/router"

func SetForumRouting(r *router.Router, h *Handlers) {
	r.GET("/api/forum/{slug}/details", h.GetForum)
	r.GET("/api/forum/{slug}/threads", h.GetThreads)
	r.GET("/api/forum/{slug}/users", h.GetForumUsers)

	r.POST("/api/forum/create", h.CreateNewForum)
	r.POST("/api/forum/{slug}/create", h.CreateNewThread)
}
