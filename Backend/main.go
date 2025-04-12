// backend/main.go
package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"
)


func main() {
	r := chi.NewRouter()

	err := godotenv.Load()
	if err!=nil {
		log.Fatal("error loading .env file")
	}

	dburl := os.Getenv("DB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx,dburl)
	if err!=nil {
		log.Fatalf("unable to connect to database %v\n",err)
	}
	defer conn.Close(ctx)

	var greeting string
    err = conn.QueryRow(ctx, "SELECT 'Hello from PostgreSQL'").Scan(&greeting)

	fmt.Println(greeting)

	port := os.Getenv("PORT");


	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(greeting))
	})

	fmt.Println("Server running on http://localhost:"+port)
	http.ListenAndServe(":"+port, r)
}
