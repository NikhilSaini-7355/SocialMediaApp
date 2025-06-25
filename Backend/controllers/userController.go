package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strconv"

	"database/sql"

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

func UpdateUser(w http.ResponseWriter, r *http.Request){

	paramID := chi.URLParam(r, "id")
	userID := r.Context().Value(middlewares.UserIDKey).(int)
	userIDstr := strconv.Itoa(userID)

	if paramID != userIDstr {
		http.Error(w,"you cannot update other's profile",http.StatusInternalServerError)
		return
	}

	var reqbody struct{
		Username string `json:"username"`
		Name string    `json:"name"`
		Email string   `json:"email"`
		Password string `json:"password"`
		Profile_pic string `json:"profile_pic"`
		Bio string `json:"bio"`
	}
 
	type UserResponse struct{
		ID int `json:"id"`
		Username string `json:"username"`
		Name string    `json:"name"`
		Email string   `json:"email"`
		Password string `json:"password"`
		Profile_pic string `json:"profile_pic"`
		Bio string `json:"bio"`
	}
	
	err := json.NewDecoder(r.Body).Decode(&reqbody)

	if err != nil {
		http.Error(w,"Invalid Json", http.StatusBadRequest)
		return
	}

	query := `SELECT name, username, email, profile_pic, bio FROM users WHERE id = $1`
	var username string
	var name string
	var email string
	var profilepic sql.NullString
	var bio sql.NullString
	err2 := database.Conn.QueryRow(context.Background(),query,userIDstr).Scan(&name, &username, &email, &profilepic, &bio)
	if err2 != nil {
	fmt.Println("Actual DB error:", err2) // 👈 print actual error
	http.Error(w, "Database error", http.StatusInternalServerError)
	return
	}

	if reqbody.Username!="" {
		username = reqbody.Username
	}
	if reqbody.Name!="" {
		name = reqbody.Name
	}
	if reqbody.Email!="" {
		email = reqbody.Email
	}
	if reqbody.Profile_pic!="" {
		profilepic.String = reqbody.Profile_pic
	}
	if reqbody.Bio!="" {
		bio.String = reqbody.Bio
	}

	if reqbody.Password!="" {
		var updated_id int
		var updated_password string
		var updated_username string
		var updated_name string
		var updated_email string
		var updated_profilepic sql.NullString
		var updated_bio sql.NullString

		 hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(reqbody.Password), bcrypt.DefaultCost)
		 query = `UPDATE users SET name = $1, username = $2, email = $3, profile_pic = $4, bio = $5, password = $6 WHERE id = $7 RETURNING id, name, username, email, profile_pic, bio, password`
		 err := database.Conn.QueryRow(context.Background(),query,name,username,email,profilepic.String,bio.String,hashedPassword,userID).Scan(&updated_id, &updated_name, &updated_username, &updated_email, &updated_profilepic, &updated_bio, &updated_password)
		 if err!=nil {
		http.Error(w, "Error saving data to database", http.StatusInternalServerError)
		return
	   }

	   updated_user := UserResponse{
		 ID : updated_id,
		 Username : updated_username,
		 Name : updated_name,
		 Email : updated_email,
		 Password : updated_password,
		 Profile_pic : updated_profilepic.String,
		 Bio : updated_bio.String,
		}

	   w.Header().Set("Content-Type", "application/json")
	   w.WriteHeader(http.StatusCreated)
	   json.NewEncoder(w).Encode(updated_user)
	} else {
		var updated_id int
		var updated_password string
		var updated_username string
		var updated_name string
		var updated_email string
		var updated_profilepic sql.NullString
		var updated_bio sql.NullString

		 query = `UPDATE users SET name = $1, username = $2, email = $3, profile_pic = $4, bio = $5 WHERE id = $6 RETURNING id, name, username, email, profile_pic, bio, password`
		 err := database.Conn.QueryRow(context.Background(),query,name,username,email,profilepic.String,bio.String,userID).Scan(&updated_id, &updated_name, &updated_username, &updated_email, &updated_profilepic, &updated_bio, &updated_password)
		 if err!=nil {
		http.Error(w, "Error saving data to database", http.StatusInternalServerError)
		return
	   }

	   updated_user := UserResponse{
		 ID : updated_id,
		 Username : updated_username,
		 Name : updated_name,
		 Email : updated_email,
		 Password : updated_password,
		 Profile_pic : updated_profilepic.String,
		 Bio : updated_bio.String,
	   }
 
	   w.Header().Set("Content-Type", "application/json")
	   w.WriteHeader(http.StatusCreated)
	   json.NewEncoder(w).Encode(updated_user)
	}
}

func GetUserProfile(w http.ResponseWriter, r *http.Request){
	username := chi.URLParam(r, "username")

	query := `SELECT id, name, username, email, profile_pic, followers, following, bio FROM users WHERE username = $1`

	var user struct{
		ID int `json:"id"`
		Name string `json:"name"`
		Username string `json:"username"`
		Email string `json:"email"`
		Profilepic sql.NullString `json:"profile_pic"`
		Followers []string `json:"followers"`
		Following []string `json:"following"`
		Bio sql.NullString `json:"bio"`
	}

	type UserResponse struct{
		ID int `json:"id"`
		Name string `json:"name"`
		Username string `json:"username"`
		Email string `json:"email"`
		Profilepic string `json:"profile_pic"`
		Followers []string `json:"followers"`
		Following []string `json:"following"`
		Bio string `json:"bio"`
	}
	err := database.Conn.QueryRow(context.Background(),query,username).Scan(&user.ID,&user.Name,&user.Username,&user.Email,&user.Profilepic,&user.Followers,&user.Following,&user.Bio)
	if err!=nil {
		fmt.Println("error: ",err)
		http.Error(w, "User Profile not found", http.StatusBadRequest)
		return
	}

	foundUser := UserResponse{
		ID : user.ID,
		Name : user.Name,
		Username : user.Username,
		Email : user.Email,
		Profilepic : user.Profilepic.String,
		Followers : user.Followers,
		Following : user.Following,
		Bio : user.Bio.String,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(foundUser)
}