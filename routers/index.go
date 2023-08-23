package routers

import (
	"net/http"

	"github.com/anirudhgray/balkan-assignment/controllers"
	"github.com/anirudhgray/balkan-assignment/routers/middleware"
	"github.com/gin-gonic/gin"
)

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
		auth.GET("/verify", controllers.VerifyEmail)
	}

	admin := v1.Group("/admin")
	{
		admin.Use(middleware.AdminAuthMiddleware())
		admin.GET("/", controllers.GetData) // TODO implement some sort of admin dashboard statistics thing here?
		admin.POST("/books", controllers.CreateBook)
		admin.PATCH("/books/:id", controllers.EditBook)
	}

	books := v1.Group("/books")
	{
		books.Use(middleware.BaseAuthMiddleware())
		books.GET("/catalog", controllers.GetBooks)
		books.POST("/cart/:bookID", controllers.AddBookToCart)
		books.GET("/cart", controllers.GetCart)
	}

	//Add All route
	//TestRoutes(route)
}
