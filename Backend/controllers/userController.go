package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strconv"

	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/database"
	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/middlewares"
	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/utils/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
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

func Login(w http.ResponseWriter, r *http.Request){
	var reqbody struct{
		Username string `json:"username"`
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

	if reqbody.Password == "" || reqbody.Username == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	query := `SELECT id, email, name, password FROM users WHERE username = $1`
	var userID int
	var hashedPassword string
	var name string
	var email string

	err = database.Conn.QueryRow(context.Background(), query, reqbody.Username).Scan(&userID, &email, &name, &hashedPassword)
	if err != nil {
		http.Error(w, "username or password wrong", http.StatusConflict)
		return
	}

	err_password := bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(reqbody.Password))
	if err_password!=nil{
		http.Error(w, "username or password wrong", http.StatusConflict)
		return
	}

	token, err := helpers.GenerateTokenAndCookie(userID,reqbody.Username)
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
		ID:       userID,
		Email:    email,
		Username: reqbody.Username,
		Name : name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func Logout(w http.ResponseWriter, r *http.Request){
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Logged out successfully",
	})
}

func FollowUnfollowUser(w http.ResponseWriter, r *http.Request){
	// followers and following arrays have id's in string format 
	// convert to int if required 
	// here, while passing the integer ids , they are first converted to string(userID) or they themselves came as string(id2)
	id2 := chi.URLParam(r, "id")
	userIDint := r.Context().Value(middlewares.UserIDKey).(int)
	userID := strconv.Itoa(userIDint)

	if id2 == userID {
		http.Error(w,"you cannot follow/unfollow yourself",http.StatusInternalServerError)
		return
	}

	var alreadyFollowing bool 
	// Check if followerID already in targetUserID's followers
	fmt.Println("FollowerID:", id2)
    fmt.Println("UserID:", userID)
	query := `SELECT $1 = ANY(COALESCE(following, '{}')) FROM users WHERE id = $2`
	err3 := database.Conn.QueryRow(context.Background(), query, id2, userID).Scan(&alreadyFollowing)
	if err3 != nil {
		  http.Error(w,err3.Error(),http.StatusInternalServerError)
		  return
	}

	if alreadyFollowing {
		_,err := database.Conn.Exec(context.Background(),`UPDATE users SET followers = array_remove(followers, $1) WHERE id = $2`,userID, id2)
		if err!= nil {
			http.Error(w,"Error unfollowing user",http.StatusInternalServerError)
			return
		}

		_,err2 := database.Conn.Exec(context.Background(),`UPDATE users SET following = array_remove(following, $1) WHERE id = $2`, id2, userID)
		if err2!= nil {
			http.Error(w,"Error unfollowing user",http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Unfollowed user"})
	} else {


		_, err := database.Conn.Exec(context.Background(),
			`UPDATE users SET following = array_append(following, $1) WHERE id = $2`, id2, userID)
		if err != nil {
			http.Error(w, "Error following user2", http.StatusInternalServerError)
			return
		}

		_, err2 := database.Conn.Exec(context.Background(),
			`UPDATE users SET followers = array_append(followers, $1) WHERE id = $2`, userID, id2)
		if err2 != nil {
			http.Error(w, "Error following user1", http.StatusInternalServerError)
			return
		}

		

		json.NewEncoder(w).Encode(map[string]string{"message": "Followed user"})
	}
}