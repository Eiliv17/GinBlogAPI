package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

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

	// get user info from context
	rawvalue, exist := c.Get("user")
	if !exist {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "not authorized",
		})
		return
	}

	user := rawvalue.(models.User)

	// body data struct
	body := struct {
		PostTitle        string   `json:"postTitle" binding:"required"`
		Tags             []string `json:"tags" binding:"required"`
		ShortDescription string   `json:"shortDescription" binding:"required"`
		PostData         string   `json:"postData" binding:"required"`
	}{}

	// gets the body data
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
			"postDocument": gin.H{
				"postTitle":        "string",
				"tags":             "[]string",
				"shortDescription": "string",
				"postData":         "string",
			},
		})
		return
	}

	// create the post object
	timeNow := time.Now()
	post := models.Post{
		PostID:           primitive.NewObjectIDFromTimestamp(timeNow),
		PostTitle:        body.PostTitle,
		Date:             timeNow,
		AuthorID:         user.UserID,
		Tags:             body.Tags,
		ShortDescription: body.ShortDescription,
		PostData:         body.PostData,
	}

	// store the post into the database
	_, err = coll.InsertOne(context.TODO(), post)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "could not store into the database",
		})
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"result": "post created",
	})
}

func DeletePost(c *gin.Context) {
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

	// get user info from context
	rawvalue, exist := c.Get("user")
	if !exist {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "not authorized",
		})
		return
	}

	user := rawvalue.(models.User)

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

	// check for post author
	if post.AuthorID != user.UserID {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "you're not the author of the post",
		})
		return
	}

	// delete post
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "post not found",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"result": "post deleted",
	})
}
