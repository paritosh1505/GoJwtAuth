package main

import (
	"jwtauth/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Output": "200 code found"})
	}

}
func main() {
	routeVal := gin.New()
	routeVal.Use(gin.Logger())
	routes.AuthSceanrio(routeVal)
	routes.RouteUser(routeVal)
	routeVal.Run(":8080")

}
