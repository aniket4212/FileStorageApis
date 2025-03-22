package controller

import (
	"filestorage/config"
	"filestorage/db/MySql"
	"filestorage/model"
	"filestorage/utils/password"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUserHandler(c *gin.Context) {
	var payload model.UserReqBody

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	hashedPass, err := password.HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	err = MySql.RegisterUserIfNotExists(payload.Username, hashedPass, config.AppConfig.DefaultStorageQuota)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User inserted or already exists"})
}
