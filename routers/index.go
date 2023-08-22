package routers

import (
	"net/http"

	"github.com/anirudhgray/balkan-assignment/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })

	v1 := route.Group("/v1")

	example := v1.Group("/example")
	{
		example.GET("/", controllers.GetData)
		example.POST("/", controllers.Create)
		example.GET("/:pid", controllers.GetSingleData)
		example.PATCH("/:pid", controllers.Update)
	}

	auth := v1.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	//Add All route
	//TestRoutes(route)
}
