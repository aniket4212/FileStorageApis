package controller

import (
	"filestorage/db/MySql"
	"filestorage/utils/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStorageHandler(c *gin.Context) {

	username := c.GetString("userName")
	if username == "" {
		fmt.Println("Username not found in token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Unauthorized"})
		return
	}
	dbData, err := MySql.FetchUserDetailsFromDB(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve storage details"})
		return
	}

	totalStorageBytes := dbData.StorageQuota
	usedStorageBytes := dbData.UsedStorage

	remainingStorageBytes := totalStorageBytes - usedStorageBytes
	totalStorageInMB := services.ConvertToAppropriateUnit(dbData.StorageQuota)

	remainingStorageInMB := services.ConvertToAppropriateUnit(remainingStorageBytes)

	c.JSON(http.StatusOK, gin.H{
		"total_storage":     totalStorageInMB,
		"remaining_storage": remainingStorageInMB,
	})
}
