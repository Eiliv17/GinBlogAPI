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
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("users")

	// get the email, username and password from req body
	body := struct {
		Email    string `json:"email" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	err := c.BindJSON(&body)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})

		return
	}

	// check if user or email already exist in database
	emailFilter := bson.D{primitive.E{Key: "email", Value: body.Email}}
	usernameFilter := bson.D{primitive.E{Key: "username", Value: body.Username}}

	// check for username
	usercount, _ := coll.CountDocuments(context.TODO(), usernameFilter)
	if usercount > 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "user already registered",
		})
		return
	}

	// check for email
	emailcount, _ := coll.CountDocuments(context.TODO(), emailFilter)
	if emailcount > 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "email already registered",
		})
		return
	}

	// hash the password using bcrypt
	hashPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "hashing password failed",
		})
		return
	}

	// create the user
	timeNow := time.Now()
	user := models.User{
		UserID:    primitive.NewObjectIDFromTimestamp(timeNow),
		Username:  body.Username,
		Email:     body.Email,
		Password:  string(hashPass),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	// stores the user inside the database
	_, err = coll.InsertOne(context.TODO(), user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "failed creating user",
		})
		return
	}

	// respond
	c.IndentedJSON(http.StatusOK, gin.H{
		"result": "user created successfully",
	})
}
