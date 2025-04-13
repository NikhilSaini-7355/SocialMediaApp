package middlewares

import (
	"context"
	"net/http"
	"fmt"
	"os"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)
        
type contextKey string
const UserIDKey = contextKey("userID")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		err := godotenv.Load()
		if err!=nil {
			log.Fatal("error loading .env file")
		}

		jwtSecret := []byte(os.Getenv("SECRET_KEY"))

		cookie, err := r.Cookie("token")
		if err!=nil {
			http.Error(w, "token not found", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		fmt.Println(tokenStr)
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		fmt.Println(token.Valid)
		if err!=nil || !token.Valid {
			http.Error(w,"Invalid Token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w,"Invalid token claims",http.StatusUnauthorized)
			return
		}

		userID := int(claims["user_id"].(float64))

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
