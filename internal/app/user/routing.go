package user

import "github.com/fasthttp/router"

func SetUserRouting(r *router.Router, h *Handlers) {
	r.POST("/api/user/{nickname}/create", h.CreateNewUser)
	r.GET("/api/user/{nickname}/profile", h.GetUserProfile)
	r.POST("/api/user/{nickname}/profile", h.UpdateUserProfile)
}
