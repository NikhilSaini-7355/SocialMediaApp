package routes

import (
	// "encoding/json"
	// "fmt"
	"net/http"

	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/middlewares"
)

func PostRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/{id}",controllers.GetPost)

	r.Group(func(r chi.Router){
		r.Use(middlewares.AuthMiddleware)
		r.Get("/feed",controllers.GetFeedPosts)
	    r.Post("/create", controllers.CreatePost)
		r.Post("/like/{id}", controllers.LikeUnlikePost)
		r.Post("/reply/{id}",controllers.ReplyToPost)
		r.Delete("/{id}",controllers.DeletePost)
	})
	return r
}
