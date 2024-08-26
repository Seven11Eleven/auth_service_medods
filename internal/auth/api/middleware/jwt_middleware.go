package middleware

import (
	"net/http"
	"strings"

	"github.com/Seven11Eleven/auth_service_medods/internal/utils"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(jwtUtils utils.JWTUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHead := c.Request.Header.Get("Authorization")
		token := strings.Split(authHead, " ")
		if len(token) == 2 {
			authToken := token[1]
			isAuthorized, err := jwtUtils.IsAuthorized(authToken)
			if isAuthorized {
				userID, err := jwtUtils.ExtractIDFromToken(authToken)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
					c.Abort()
					return
				}
				c.Set("userID", userID.String())
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"errpr": err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unathorized"})
		c.Abort()
	}
}
