package middleware

import (
	"filestorage/utils/generateToken"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticateForToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		log.Println("Token is empty")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Unauthorized"})
		return
	}

	// Fetching data from token
	userName, shortTokenVerificationErr := generateToken.VerifyShortToken(token)

	if shortTokenVerificationErr != nil {
		log.Println("Token verification failed :::", shortTokenVerificationErr)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Unauthorized"})
		return
	}

	c.Set("userName", userName)
	c.Next()
}
