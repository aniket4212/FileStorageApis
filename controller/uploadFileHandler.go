package controller

import (
	"filestorage/config"
	"filestorage/db/MySql"
	"filestorage/model"
	"filestorage/utils/services"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func FileUploadHandler(c *gin.Context) {
	username := c.GetString("userName")
	if username == "" {
		log.Println("Username not found in token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Unauthorized"})
		return
	}

	// Retrieve uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("Error retrieving file:", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "No file uploaded"})
		return
	}
	defer file.Close()

	fileSize := header.Size
	originalFileName := header.Filename

	// Fetch users used storage
	dbStorage, err := MySql.FetchUserDetailsFromDB(username)
	if err != nil {
		log.Println("Error fetching user storage:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Storage validation error"})
		return
	}

	remaining_storage := dbStorage.StorageQuota - dbStorage.UsedStorage

	// Check if enough storage is available for upload file
	if fileSize > remaining_storage {
		fmt.Println("Not enough storage space")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "message": "Not enough storage space"})
		return
	}

	// Create User Directory If Not Exists
	currentDir, err := os.Getwd()
	fmt.Println("currentDir::", currentDir)
	if err != nil {
		log.Println("Error fetching user home directory:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Server error"})
		return
	}

	baseDir := config.AppConfig.BaseDir

	userDir := filepath.Join(currentDir, baseDir, username)

	// Ensure that users directory exist if doesn't exist create it.
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		err = os.MkdirAll(userDir, os.ModePerm)
		if err != nil {
			log.Println("Error creating user directory:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error creating user folder"})
			return
		}
	}

	name, extension := SplitFilenameAndExtension(originalFileName)
	//Prevent file overwriting if same file uploaded
	timestamp := time.Now().Format("20060102150405")
	newFileName := fmt.Sprintf("%s_%s%s", name, timestamp, extension)
	filePath := filepath.Join(userDir, newFileName)

	// Save the file to user directory
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Error saving file:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error saving file"})
		return
	}
	defer outFile.Close()

	//Copy file contents
	_, err = io.Copy(outFile, file)
	if err != nil {
		log.Println("Error writing file:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error writing file"})
		return
	}

	err = MySql.UpdateUserStorage(username, fileSize)
	if err != nil {
		log.Println("Error updating user storage:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Storage update failed"})
		return
	}

	err = MySql.StoreFileMetadata(model.FileMetadata{
		FileName:         newFileName,
		OriginalFileName: originalFileName,
		UploadedBy:       username,
		Size:             fileSize,
		UploadTime:       time.Now().Format(time.RFC3339),
	})

	if err != nil {
		log.Println("Error storing metadata:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error storing metadata"})
		return
	}

	sizeString := services.ConvertToAppropriateUnit(fileSize)

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"message":  "File uploaded successfully",
		"filename": newFileName,
		"sizeInMB": sizeString,
	})
	fmt.Println("File uploaded successfully")
}
func SplitFilenameAndExtension(filename string) (name string, extension string) {
	// Get the file extension (including the dot, e.g., ".txt")
	extension = filepath.Ext(filename)

	// Remove the extension from the filename to get the name
	name = filename[:len(filename)-len(extension)]

	return name, extension
}
