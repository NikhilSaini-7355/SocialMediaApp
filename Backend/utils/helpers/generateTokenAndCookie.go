package helpers

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"github.com/joho/godotenv"
	"fmt"
)


func GenerateTokenAndCookie(userId int, username string)(string,error){
	err := godotenv.Load()
	if err!=nil {
		log.Fatal("error loading .env file")
	}

	jwtSecret := []byte(os.Getenv("SECRET_KEY"))
	claims := jwt.MapClaims{
		"user_id" : userId,
		"username" : username,
		"exp" : time.Now().Add(time.Hour * 360).Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	signedToken, err := token.SignedString(jwtSecret)

	fmt.Println(signedToken)
	if err!=nil {
		return "",err
	}

	return signedToken,nil
}