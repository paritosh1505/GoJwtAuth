package routes

import (
	"jwtauth/controller"
	"jwtauth/middleware"

	"github.com/gin-gonic/gin"
)

func routeUser(incomingReq *gin.Engine) {
	incomingReq.Use(middleware.MyAuthentication())
	incomingReq.GET("/usersval/:userId", controller.GetUserDetail())
	incomingReq.GET("/usersval", controller.AggreGator())
}
