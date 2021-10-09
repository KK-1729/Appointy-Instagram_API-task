package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	UID primitive.ObjectID `json:"uid,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Posts struct{
		AllPosts []Post `json:"all_posts" bson:"all_posts"`
		Count int64 `json:"count" bson:"count"`
		PageFormat struct{
			HasNextPage bool `json:"has_next_page" bson:"has_next_page"`
			EndCursor string `json:"end_cursor" bson:"end_cursor"`
		} `json:"page_format" bson:"page_format"`
	} `json:"posts" bson:"posts"`
}

type Post struct {
	PID primitive.ObjectID `json:"pid,omitempty" bson:"_id,omitempty"`
	Caption string `json:"caption,omitempty" bson:"caption,omitempty"`
	Image string `json:"image,omitempty" bson:"image,omitempty"`
	Time time.Time `json:"time,omitempty" bson:"time,omitempty"`
	//User *User `json:"user,omitempty" bson:"user,omitempty"`
}

type UserResp struct {
	UserInfo User `json:"user_info" bson:"user_info"`
}
