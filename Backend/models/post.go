package models

type Post struct{
	ID int `json:"id"`
	PostedBy int `json:"postedBy"` 
	Text string `json:"text"`
	Img  string `json:"img"`
	Likes int `json:"likes"`
	Replies []Reply `json:"replies"`
	CreatedAt string `json:"created_at"`
}