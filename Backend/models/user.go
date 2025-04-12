package models

type User struct {
	ID       int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`           // "-" hides it from JSON responses
	CreatedAt string `json:"created_at"`  // optional timestamp
	ProfilePic string `json:"profilePic"`
	Followers  []string `json:"followers"`
	Following  []string `json:"following"`
	Bio        string   `json:"bio"`
}
