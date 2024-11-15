package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"whattoday/web-service-gin/users"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Modify RequireAuth to return gin.HandlerFunc
func RequireAuth(userDao users.UserDao) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the cookie of the request
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Decode/validate it
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			// Use the secret for validation
			return []byte(os.Getenv("SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check expiration
			fmt.Println("claims")
			fmt.Println(claims)
			fmt.Println("claims")
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// Find the user with token sub
			userIdFloat := claims["sub"].(float64) // Type assert the value
			userId := int(userIdFloat)
			user, err := userDao.GetUser(userId)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// Attach to request context
			c.Set("user", user)
			c.Next()

		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
