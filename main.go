package main

import (
	"github.com/Eiliv17/GinJWTAuthAPI/controllers"
	"github.com/Eiliv17/GinJWTAuthAPI/initializers"
	"github.com/Eiliv17/GinJWTAuthAPI/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	userr := r.Group("/user")
	{
		userr.POST("/signup", controllers.Signup)

		userr.POST("/login", controllers.Login)

	}

	// authenticated test route
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	blog := r.Group("/posts")
	{
		// default blog routes
		blog.GET("", controllers.GetPosts)

		blog.GET("/:id", controllers.GetPost)

		// authenticated blog routes routes
		blog.POST("", middleware.RequireAuth, controllers.CreatePost)

		blog.DELETE("/:id", middleware.RequireAuth, controllers.DeletePost)
	}

	r.Run() // listen and serve on port specified by PORT env var
}
