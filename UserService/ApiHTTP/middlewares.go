package ApiHTTP

import (
	"../utils"
	"github.com/gin-gonic/gin"

	conf "../Config"
)



// JWTAuthMiddleware, check the jwt token string from request.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get("Auth-Token")
		if authToken == "" {
			c.JSON(400, gin.H{"error": " Auth-Token header required",})
			c.Abort()
			return
		}

		// parseToken
		claims, err := utils.ParseJWTToken(authToken, []byte(conf.AuthTokenSalt))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		// call next handler function
		c.Set(JWTGetUserId, claims.Id)
		c.Next()
	}
}
