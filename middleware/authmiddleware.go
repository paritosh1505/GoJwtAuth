package middleware

import (
	"jwtauth/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MyAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Error while getting the token"})
			return
		}
		claims, err := helper.GetClaimfromToken(clientToken)
		if err != "" {
			c.JSON(http.StatusBadGateway, gin.H{"Error": "Error while getting claim from token"})
			return
		}
		c.Set("Name", claims.Name)
		c.Set("Name", claims.Subject)
		c.Set("Name", claims.Email)
	}

}
