package main

import (
	"github.com/Eiliv17/GinJWTAuthAPI/controllers"
	"github.com/Eiliv17/GinJWTAuthAPI/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	v1 := r.Group("/user")
	{
		v1.POST("/signup", controllers.Signup)

		//v1.POST("/login", loginEndpoint)
	}

	r.Run() // listen and serve on port specified by PORT env var
}
