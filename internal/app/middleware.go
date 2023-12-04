package app

import (
	"net/http"
	"project/internal/pkg/roles"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const prefix = "Bearer"

func (a *Application) OnAuthCheck(allowedRoles ...roles.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractTokenFromHeader(c.Request)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("SuperSecretKey"), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID, ok := claims["userID"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Set("UserID", int(userID))

		userRole, ok := claims["role"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if !isRoleAllowed(int(userRole), allowedRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "доступ отказан"})
			return
		}

		c.Next()
	}
}

func extractTokenFromHeader(req *http.Request) string {
	bearerToken := req.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}

	if strings.Split(bearerToken, " ")[0] != prefix {
		return ""
	}

	return strings.Split(bearerToken, " ")[1]
}

func isRoleAllowed(userRole int, allowedRoles []roles.Role) bool {
	for _, allowedRole := range allowedRoles {
		if userRole == int(allowedRole) {
			return true
		}
	}
	return false
}
