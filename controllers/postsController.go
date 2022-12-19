package controllers

import (
	"context"
	"net/http"
	"os"

	"github.com/Eiliv17/GinJWTAuthAPI/initializers"
	"github.com/Eiliv17/GinJWTAuthAPI/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPosts(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("posts")

	// retrieve posts from database
	cursor, _ := coll.Find(context.TODO(), bson.D{})

	// retrieve all posts from cursor
	var posts []models.Post
	_ = cursor.All(context.TODO(), &posts)

	// checks if documents are returned
	if len(posts) == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "no post found",
		})
		return
	}

	// send posts
	c.IndentedJSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("posts")

	// get post id parameter
	postID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid post id",
		})
		return
	}

	// retrieve post from database
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	result := coll.FindOne(context.TODO(), filter)

	var post models.Post
	err = result.Decode(&post)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "post not found",
		})
		return
	}

	// return the post
	c.IndentedJSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("posts")

}

/* func DeletePost(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("posts")

}  */
