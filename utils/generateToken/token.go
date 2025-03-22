package generateToken

import (
	"filestorage/config"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userName": userName,
		"exp":      time.Now().Unix() + int64(config.AppConfig.TokenValidityInSeconds),
	})

	return token.SignedString([]byte(config.AppConfig.SecretKey))
}

func VerifyShortToken(token string) (string, error) {
	parsedToken, parsedTokenErr := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(config.AppConfig.SecretKey), nil
	})

	if parsedTokenErr != nil {
		log.Println("Error in parsing token =-=-=-=-=-=", parsedTokenErr)
		return "", fmt.Errorf("error parsing token")
	}

	isTokenValid := parsedToken.Valid

	if !isTokenValid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("claim token error")
	}

	userName, exists := claims["userName"].(string)
	if !exists {
		log.Println("userName key missing in token claims")
		return "", fmt.Errorf("userName not found in token")
	}

	return userName, nil
}
