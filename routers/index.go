package routers

import (
	"net/http"

	"github.com/anirudhgray/balkan-assignment/controllers"
	"github.com/anirudhgray/balkan-assignment/routers/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })

	v1 := route.Group("/v1")

	v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := v1.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/verify", controllers.VerifyEmail)
		auth.GET("/request-verification", controllers.RequestVerificationAgain)
		auth.GET("/forgot-password", controllers.ForgotPasswordRequest)
		auth.POST("/delete-account", controllers.DeleteAccount)          // otp
		auth.POST("/set-forgotten-password", controllers.SetNewPassword) // otp
		auth.POST("/request-account-deletion", middleware.BaseAuthMiddleware(), controllers.RequestDeletion)
		auth.POST("/reset-password", middleware.BaseAuthMiddleware(), controllers.ResetPassword)
	}

	admin := v1.Group("/admin")
	{
		admin.Use(middleware.AdminAuthMiddleware())
		// implement some sort of admin dashboard statistics thing here?
		admin.GET("/users", controllers.GetAllUsers)
		admin.POST("/users/:userID", controllers.BanUser)
		admin.POST("/verify/user/:id", controllers.ManualVerification)
		admin.POST("/books", controllers.CreateBook)
		admin.PATCH("/books/:id", controllers.EditBook)
		admin.GET("/transactions", controllers.GetAllTransactions)
		admin.DELETE("/reviews/:id", controllers.DeleteReview)
		admin.DELETE("/books/:id", controllers.DeleteBook)
		admin.DELETE("/books/:id/hard", controllers.DeleteBookHard) // from user libraries as well.
	}

	superadmin := v1.Group("/superadmin")
	{
		superadmin.Use(middleware.SuperAdminAuthMiddleware())
		superadmin.POST("/users/promote/:userID", controllers.PromoteUserToAdmin)
		superadmin.POST("/users/demote/:userID", controllers.DemoteUserToBase)
	}

	books := v1.Group("/books")
	{
		books.Use(middleware.BaseAuthMiddleware())
		books.GET("/recommend", controllers.GenerateRecommendations)
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
		checkout.POST("/credits", controllers.BuyCredits)
	}

	user := v1.Group("/user")
	{
		user.Use(middleware.BaseAuthMiddleware())
		user.GET("/transactions", controllers.GetUserTransactions)
		user.GET("/library", controllers.GetUserLibrary)
	}
}
