package routes

import (
	"jwtauth/controller"

	"github.com/gin-gonic/gin"
)

func AuthSceanrio(incomingReq *gin.Engine) {
	incomingReq.POST("/user/signup", controller.SingupUser)
	incomingReq.POST("/user/login", controller.LoginUser)
}
