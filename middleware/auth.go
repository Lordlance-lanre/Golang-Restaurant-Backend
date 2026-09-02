package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helpers"github.com/Lordlance-lanre/Golang-Restaurant-Backend/helpers"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Next()
		clientToken := c.Request.Header.Get("token")
		if clientToken == ""{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		claims, msg := helpers.ValidateAllTokens(clientToken)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}