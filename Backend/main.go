// backend/main.go
package main

import (
	
	"fmt"
	"net/http"
	"log"
	"os"
	"context"
	
	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/routes"
	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)


func main() {
	r := chi.NewRouter()

	err := godotenv.Load()
	if err!=nil {
		log.Fatal("error loading .env file")
	}

	database.InitDB()
	
	defer database.Conn.Close(context.Background())

	port := os.Getenv("PORT");

	r.Mount("/api/users",routes.UserRouter())

	fmt.Println("Server running on http://localhost:"+port)
	http.ListenAndServe(":"+port, r)
}
