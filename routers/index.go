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
		// admin.GET("/", controllers.GetData) // implement some sort of admin dashboard statistics thing here?
		admin.GET("/users", controllers.GetAllUsers)
		admin.POST("/books", controllers.CreateBook)
		admin.PATCH("/books/:id", controllers.EditBook)
		admin.GET("/transactions", controllers.GetAllTransactions)
		admin.DELETE("/reviews/:id", controllers.DeleteReview)
		admin.DELETE("/books/:id", controllers.DeleteBook)
		admin.DELETE("/books/:id/hard", controllers.DeleteBookHard) // from user libraries as well.
	}

	books := v1.Group("/books")
	{
		books.Use(middleware.BaseAuthMiddleware())
		books.GET("/catalog", controllers.GetBooks)
		books.GET("/catalog/:bookID", controllers.GetBook)
		books.GET("/catalog/:bookID/download", controllers.DownloadBook)
		books.POST("/cart/add/:bookID", controllers.AddBookToCart)
		books.POST("/cart/remove/:bookID", controllers.RemoveFromCart)
		books.GET("/cart", controllers.GetCart)
		books.POST("/attach", controllers.AttachCL)
		books.POST("/review", controllers.AddReview)
	}

	checkout := v1.Group("/checkout")
	{
		checkout.Use(middleware.BaseAuthMiddleware())
		checkout.POST("/cart", controllers.Checkout)
	}

	user := v1.Group("/user")
	{
		user.Use(middleware.BaseAuthMiddleware())
		user.GET("/transactions", controllers.GetUserTransactions)
		user.GET("/library", controllers.GetUserLibrary)
	}
}
