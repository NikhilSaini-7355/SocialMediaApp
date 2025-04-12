package controllers

import (
	"encoding/json"
	"net/http"
	"context"
	"fmt"
	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/database"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgx/v5"
	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/utils/helpers"
)



func SignUp(w http.ResponseWriter, r *http.Request) {

	var reqbody struct{
		Username string `json:"username"`
		Name string    `json:"name"`
		Email string   `json:"email"`
		Password string `json:"password"`
	}

	type UserResponse struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqbody)

	if err != nil {
		http.Error(w,"Invalid Json", http.StatusBadRequest)
		return
	}

	if reqbody.Email == "" || reqbody.Password == "" || reqbody.Username == "" || reqbody.Name=="" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	query := `SELECT id FROM users WHERE email = $1 OR username = $2`
	var existingID int

	err = database.Conn.QueryRow(context.Background(), query, reqbody.Email, reqbody.Username).Scan(&existingID)
	if err == nil {
		http.Error(w, "User with this email or username already exists", http.StatusConflict)
		return
	}

	// if err != pgx.ErrNoRows, it's an actual DB error
	if  err != pgx.ErrNoRows {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(reqbody.Password), bcrypt.DefaultCost)

	var userId int
	insertQuery := `INSERT INTO users (name, email, username, password) VALUES ($1, $2, $3, $4) RETURNING id;`
    err = database.Conn.QueryRow(context.Background(), insertQuery,reqbody.Name, reqbody.Email, reqbody.Username, hashedPassword).Scan(&userId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := helpers.GenerateTokenAndCookie(userId,reqbody.Username)
	if err!=nil {
		http.Error(w,"Error generating token",http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: token,
		Path: "/",
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteLaxMode,
		MaxAge: 3600 * 24 * 15,
	})

	user := UserResponse{
		ID:       userId,
		Email:    reqbody.Email,
		Username: reqbody.Username,
		Name : reqbody.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}