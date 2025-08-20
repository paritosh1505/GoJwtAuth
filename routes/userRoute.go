package routes

import (
	"jwtauth/controller"
	"jwtauth/middleware"

	"github.com/gin-gonic/gin"
)

func RouteUser(incomingReq *gin.Engine) {
	incomingReq.Use(middleware.MyAuthentication())
	incomingReq.GET("/userval/:mongouid", controller.GetUserDetail())
	incomingReq.GET("/userval", controller.AggreGator())
}
