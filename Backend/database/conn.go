package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var Conn *pgx.Conn

func InitDB() {
	var err2 error
    err2 = godotenv.Load()
	if err2!=nil {
		log.Fatal("error loading .env file")
	}

	dburl := os.Getenv("DB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var err error
	Conn, err = pgx.Connect(ctx,dburl)   // do pgx.pool when the web app is in production 

	if err!=nil {
		log.Fatalf("unable to connect to database %v\n",err)
	}


    err = Conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}

	fmt.Println("✅ PostgreSQL connected and ping successful")
}