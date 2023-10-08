package jwttoken

import (
	"fmt"
	"net/http"
	"project/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GenerateJWTToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID

	tokenString, err := token.SignedString([]byte("SuperSecretKey"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckJWTToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("entered")
		cookie, err := c.Cookie("jwtToken")
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		fmt.Println(cookie)
		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неправильный метод подписи")
			}
			return []byte("SuperSecretKey"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, middleware.Response{
				Status:  "Failed",
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, middleware.Response{
				Status:  "Failed",
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		userID := uint(claims["userID"].(float64))

		c.Set("userID", userID)

		c.Next()
	}
}

func ParseJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("SuperSecretKey"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("неверный JWT токен")
	}
}

func GetUserIDbyToken(tokenString string) (uint, error) {
	tokenClaims, err := ParseJWTToken(tokenString)
	if err != nil {
		return 0, err
	}

	userID := tokenClaims["userID"].(float64)

	return uint(userID), nil
}
