package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Eiliv17/GinJWTAuthAPI/initializers"
	"github.com/Eiliv17/GinJWTAuthAPI/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// signup controller
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

// login controller
func Login(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("users")

	// get the email and pass off req body
	body := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	err := c.BindJSON(&body)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	// look up requested user
	emailFilter := bson.D{primitive.E{Key: "email", Value: body.Email}}
	result := coll.FindOne(context.TODO(), emailFilter)

	// decode result
	var user models.User
	err = result.Decode(&user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// compare sent pass with saved user pass hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.UserID.Hex(),
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	HMACSecret := os.Getenv("HMAC_SECRET")

	tokenString, err := token.SignedString([]byte(HMACSecret))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}

	// set token as cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	// send the JWT token
	c.IndentedJSON(http.StatusOK, gin.H{
		"result": "logged in successfully",
	})
}

// validation
func Validate(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "I'm logged in",
	})
}
