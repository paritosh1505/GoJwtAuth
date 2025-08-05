package routes

import (
	"jwtauth/controller"
	"jwtauth/middleware"

	"github.com/gin-gonic/gin"
)

func routeUser(incomingReq *gin.Engine) {
	incomingReq.Use(middleware.Authetication())
	incomingReq.GET("/usersval", controller.GetUserDetail())
	incomingReq.GET("/usersval/:userId", controller.GetUserDetail())
}
