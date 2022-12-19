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

	blog := r.Group("")
	{
		// authenticated test route
		blog.GET("/validate", middleware.RequireAuth, controllers.Validate)

		// default blog routes
		blog.GET("/posts", controllers.GetPosts)

		blog.GET("/posts/:id", controllers.GetPost)

		// authenticated blog routes routes
		/* blog.POST("/posts", middleware.RequireAuth, controllers.CreatePost) */

		/* blog.DELETE("/posts", middleware.RequireAuth, controllers.DeletePost)  */
	}

	r.Run() // listen and serve on port specified by PORT env var
}
