package models

type Reply struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	UserProfilePic string `json:"userProfilePic"`
	CreatedAt string `json:"created_at"`
}
