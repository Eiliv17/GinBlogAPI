package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	PostID           primitive.ObjectID `bson:"_id" json:"postID"`
	PostTitle        string             `bson:"postTitle" json:"postTitle"`
	Date             time.Time          `bson:"date" json:"date"`
	AuthorID         primitive.ObjectID `bson:"authorID" json:"authorID"`
	Tags             []string           `bson:"tags" json:"tags"`
	ShortDescription string             `bson:"shortDescription" json:"shortDescription"`
	PostData         string             `bson:"postData" json:"postData"`
}
