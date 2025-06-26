package controllers

import (
	"unicode/utf8"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	// "time"
	"strconv"

	"database/sql"

	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/database"
	"github.com/NikhilSaini-7355/SocialMediaApp/Backend/middlewares"
	// "github.com/NikhilSaini-7355/SocialMediaApp/Backend/utils/helpers"
	"github.com/go-chi/chi/v5"
	// "github.com/jackc/pgx/v5"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var reqbody struct{
		PostedBy int `json:"postedBy"`
		Text string `json:"text"`
		Img string `json:"img"`
	}
	
	type postResponse struct{
		ID int `json:"id"`
		PostedBy int `json:"postedBy"`
		Text string `json:"text"`
	    Img  string `json:"img"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqbody)

	if err != nil {
		http.Error(w,"Invalid Json", http.StatusBadRequest)
		return
	}

	Posted_By := strconv.Itoa(reqbody.PostedBy)

	if Posted_By=="" {
		http.Error(w,"PostedBy is required", http.StatusBadRequest)
		return
	}

	if reqbody.Text=="" {
		http.Error(w,"Text is required", http.StatusBadRequest)
		return
	}

	var Posted_User_Id int
	query := `SELECT id FROM users WHERE id = $1`
	err = database.Conn.QueryRow(context.Background(),query,reqbody.PostedBy).Scan(&Posted_User_Id)
	if err!=nil {
		http.Error(w,"cannot find posted id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middlewares.UserIDKey).(int)
	// userIDstr := strconv.Itoa(userID)

	if Posted_User_Id!=userID {
		http.Error(w,"You cannot send post from other id", http.StatusBadRequest)
		return
	}

	maxlength := 100
	textlength := utf8.RuneCountInString(reqbody.Text)

	if maxlength<textlength {
		http.Error(w,"Max Length Exceeded", http.StatusBadRequest)
		return
	}

	var newPost postResponse
	query = `INSERT INTO posts (posted_by, text, img) VALUES ($1, $2, $3) RETURNING id, posted_by, text, img`
	err = database.Conn.QueryRow(context.Background(),query,Posted_User_Id,reqbody.Text,reqbody.Img).Scan(&newPost.ID,&newPost.PostedBy,&newPost.Text,&newPost.Img)
	if err!=nil {
		fmt.Println("error: ",err)
		http.Error(w,"Unable to save post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	type postResponse struct{
		ID int `json:"id"`
		PostedBy int `json:"posted_by"`
		Text string `json:"text"`
	    Img  string `json:"img"`
		Likes []int `json:"likes"`
	}

	query := `SELECT id, posted_by, text, img, likes FROM posts WHERE id = $1`

	var post postResponse
	err := database.Conn.QueryRow(context.Background(),query,postID).Scan(&post.ID,&post.PostedBy,&post.Text,&post.Img,&post.Likes)
	if err!=nil {
		http.Error(w,"Did not found post", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	query := `SELECT id, posted_by FROM posts WHERE id = $1`

	var post_ID int
	var postedBy int
	err := database.Conn.QueryRow(context.Background(),query,postID).Scan(&post_ID,&postedBy)
	if err!=nil {
		http.Error(w,"Post not found", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middlewares.UserIDKey).(int)
	// userIDstr := strconv.Itoa(userID)
	if postedBy!=userID {
		http.Error(w,"Unauthorised to delete post", http.StatusBadRequest)
		return
	}

	query = `DELETE FROM posts WHERE id = $1`
	_,err = database.Conn.Exec(context.Background(),query,post_ID)
	if err!=nil {
		http.Error(w,"Unable to delete post", http.StatusInternalServerError)
		return
	}

	type FinalMessage struct{
		Message string `json:"message"`
	}

	message := FinalMessage{
		Message : "Post deleted successfully",
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func LikeUnlikePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)
	// userIDstr := strconv.Itoa(userID)
	postID := chi.URLParam(r, "id")

	var post_id int
	query := `SELECT id FROM posts WHERE id = $1`
	err := database.Conn.QueryRow(context.Background(),query,postID).Scan(&post_id)
	if err!=nil{
		http.Error(w,"Invalid Post ID", http.StatusBadRequest)
		return
	}

	var alreadyLiked bool 
	query = `SELECT $1 = ANY(likes) FROM posts WHERE id = $2`
	err = database.Conn.QueryRow(context.Background(),query,userID,post_id).Scan(&alreadyLiked)
	if err!=nil {
		http.Error(w,"Error knowing current like status", http.StatusInternalServerError)
		return
	}

	type FinalMessage struct{
		Message string `json:"message"`
	}

	if alreadyLiked {
		query = `UPDATE posts SET likes = array_remove(likes, $1) WHERE id = $2`
		_,err = database.Conn.Exec(context.Background(),query,userID,post_id)
		if err!=nil {
			http.Error(w,"Unable to unlike post", http.StatusBadRequest)
		    return
		}
		message := FinalMessage{
			Message : "Post Unliked",
		}
		w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusCreated)
	    json.NewEncoder(w).Encode(message)
	} else {
		query = `UPDATE posts SET likes = array_append(likes, $1) WHERE id = $2`
		_,err = database.Conn.Exec(context.Background(),query,userID,post_id)
		if err!=nil {
			http.Error(w,"Unable to like post", http.StatusBadRequest)
		    return
		}
		message := FinalMessage{
			Message : "Post Liked",
		}
		w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusCreated)
	    json.NewEncoder(w).Encode(message)
	}
}

func ReplyToPost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	userID := r.Context().Value(middlewares.UserIDKey).(int)
	post_ID,_ := strconv.Atoi(postID)

	var reqbody struct{
		Text string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqbody)
	if err != nil {
		http.Error(w,"Invalid Json", http.StatusBadRequest)
		return
	}

	var userid int
	var username string 
	var profilepic sql.NullString

	query := `SELECT id, username, profile_pic FROM users WHERE id = $1`
	err = database.Conn.QueryRow(context.Background(),query,userID).Scan(&userid,&username,&profilepic)
	if err!=nil {
		fmt.Println("error : ",err)
		http.Error(w,"Error in fetching user details", http.StatusBadRequest)
		return
	}

	type Reply struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	UserProfilePic string `json:"userProfilePic"`
   }

   var userReply Reply

   query = `INSERT INTO replies (post_id, user_id, username, text, user_profile_pic) VALUES ($1, $2, $3, $4, $5) RETURNING id, post_id, user_id, username, text, user_profile_pic`
   err = database.Conn.QueryRow(context.Background(),query,post_ID,userid,username,reqbody.Text,profilepic.String).Scan(&userReply.ID,&userReply.PostID,&userReply.UserID,&userReply.Username,&userReply.Text,&userReply.UserProfilePic)
   if err!=nil {
	 http.Error(w,"Error in replying to the post", http.StatusInternalServerError)
		return
   }

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userReply)
}

func GetFeedPosts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	query := `SELECT * FROM posts WHERE posted_by::text = ANY (SELECT unnest(following) FROM users WHERE id = $1) ORDER BY created_at DESC`

	// BETTER SOLN IS TO CONVERT THE FOLLOWER AND FOLLOWING FROM []STRING TO []INTEGER
	// THE ABOVE QUERY IS JUST A FIX
	// WILL PERFORM THIS MAJOR FIX IN FREE TIME

	rows,err := database.Conn.Query(context.Background(),query,userID)
	if err!=nil {
		fmt.Println("error: ",err)
		http.Error(w, "Error fetching feed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Post struct{
	ID int `json:"id"`
	PostedBy int `json:"posted_by"` 
	Text string `json:"text"`
	Img  sql.NullString `json:"img"`
	Likes []int `json:"likes"`
	CreatedAt string `json:"created_at"`
    }

	var posts []Post

	for rows.Next() {
		var p Post
		var tm time.Time
		err = rows.Scan(&p.ID,&p.PostedBy,&p.Text,&p.Img,&tm,&p.Likes)
		if err != nil {
			fmt.Println("error: ",err)
			http.Error(w, "Error scanning posts", http.StatusInternalServerError)
			return
		}
		p.CreatedAt = tm.Format(time.RFC3339)
		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(posts)
}