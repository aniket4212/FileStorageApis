package controller

import (
	"filestorage/db/MySql"
	"filestorage/model"
	"filestorage/utils/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUploadedFilesHandler(c *gin.Context) {
	username := c.GetString("userName")
	if username == "" {
		fmt.Println("Username not found in token")
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Unauthorized"})
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	sqlOffset := offset - 1
	if sqlOffset < 0 {
		sqlOffset = 0
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	files, totalRecords, err := MySql.FetchUploadedFiles(username, limit, sqlOffset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to retrieve files"})
		return
	}

	var responseFiles []model.FileResponse
	for _, file := range files {
		responseFiles = append(responseFiles, model.FileResponse{
			FileName:   file.OriginalFileName,
			UploadedBy: file.UploadedBy,
			Size:       services.ConvertToAppropriateUnit(file.Size),
			UploadTime: file.UploadTime,
		})
	}

	totalPages := (totalRecords + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"status":      true,
		"files":       responseFiles,
		"offset":      sqlOffset,
		"limit":       limit,
		"total_files": totalRecords,
		"total_pages": totalPages,
	})
}
