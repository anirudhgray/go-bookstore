package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func GenerateToken(user models.User) (string, error) {

	tokenLifespan, err := strconv.Atoi(viper.GetString("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(viper.GetString("API_SECRET")))

}

// TODO auth middleware for routes
func getTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func GetToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(viper.GetString("API_SECRET")), nil
	})
	return token, err
}

func ValidateToken(c *gin.Context) (uint, error) {
	token, err := GetToken(c)

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id, _ := claims["id"].(float64) //get the id
		return uint(id), nil
	}

	return 0, errors.New("invalid token provided")
}
