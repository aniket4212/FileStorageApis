package controller

import (
	"filestorage/db/MySql"
	"filestorage/model"
	"filestorage/utils/generateToken"
	"filestorage/utils/password"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GenerateToken(c *gin.Context) {
	var RequestBodyForToken model.UserReqBody

	err := c.ShouldBindJSON(&RequestBodyForToken)
	if err != nil {
		log.Println("Error in request data payload =-=-=-=", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid Data Received"})
		return
	}

	userDetailsFromDb, err := MySql.FetchUserDetailsFromDB(RequestBodyForToken.Username)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "No Username found"})
			return
		}
		log.Println("Error fetching user details:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Database error"})
		return
	}

	checkPassword := password.CompareHash(RequestBodyForToken.Password, userDetailsFromDb.Password)
	if !checkPassword {
		log.Println("Invalid password for user:", RequestBodyForToken.Username)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid Password"})
		return
	}

	token, err := generateToken.GenerateToken(RequestBodyForToken.Username)
	if err != nil {
		log.Println("Error generating token:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "token": token})
}
