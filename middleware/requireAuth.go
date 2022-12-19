package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Eiliv17/GinJWTAuthAPI/initializers"
	"github.com/Eiliv17/GinJWTAuthAPI/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequireAuth(c *gin.Context) {
	// database setup
	dbname := os.Getenv("DB_NAME")
	coll := initializers.DB.Database(dbname).Collection("users")

	// Get the cookie off request headers
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "not authorized",
		})
		return
	}

	// decode/validate it
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("HMAC_SECRET")), nil
	})

	if token == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "not authorized",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check the expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "not authorized",
			})
			return
		}

		// find the user with token sub
		var user models.User
		returnedID, _ := primitive.ObjectIDFromHex(claims["userID"].(string))
		userIDFileter := bson.D{primitive.E{Key: "_id", Value: returnedID}}
		result := coll.FindOne(context.TODO(), userIDFileter)
		err := result.Decode(&user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "not authorized",
			})
			return
		}

		// attach user to the context
		c.Set("user", user)

		// continue
		c.Next()

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "not authorized",
		})
		return
	}

	c.Next()
}
