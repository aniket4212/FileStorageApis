package routes

import (
	"filestorage/config"
	"filestorage/controller"
	"filestorage/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(server *gin.Engine) {
	// Ping API
	server.GET(config.AppConfig.Prefix+"/ping", PingResponse)

	server.POST(config.AppConfig.Prefix+"/register", controller.RegisterUserHandler)
	server.POST(config.AppConfig.Prefix+"/login", controller.GenerateToken)

	authRoutes := server.Group(config.AppConfig.Prefix)
	authRoutes.Use(middleware.AuthenticateForToken)
	{
		authRoutes.POST("/upload", controller.FileUploadHandler)
		authRoutes.GET("/storage/remaining", controller.GetStorageHandler)
		authRoutes.GET("/files", controller.GetUploadedFilesHandler)

	}

}

// Function to Check PING Response
func PingResponse(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "PONG"})
}
