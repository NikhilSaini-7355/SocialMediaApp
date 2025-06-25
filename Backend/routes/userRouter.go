package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/middlewares"
)

func UserRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/profile/{username}",controllers.GetUserProfile)

	r.Post("/signup", controllers.SignUp)
	r.Post("/login", controllers.Login)
	r.Post("/logout", controllers.Logout)

	r.Group(func(r chi.Router){
		r.Use(middlewares.AuthMiddleware)
	    r.Post("/follow/{id}", controllers.FollowUnfollowUser)
		r.Post("/update/{id}", controllers.UpdateUser)
	})


	type reqbody2 struct{
		Message string `json:"message"`
	}
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit /api/users/")
		message2 := reqbody2{
			Message : "hello",
		}
		w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusOK)
	    json.NewEncoder(w).Encode(message2)
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.Write([]byte(fmt.Sprintf("Get user with ID: %s", id)))
	})

	return r
}
