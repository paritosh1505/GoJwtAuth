package main

import (
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
	routeVal.GET("/", GetCheck())
	routeVal.Run(":8080")

}
