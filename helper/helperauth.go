package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func MatchUserId(c *gin.Context, userid string) error {
	userType := c.GetString("userType")
	uid := c.GetString("userid")
	var errorval error = nil
	if userType == "USER" && uid != userid {
		errorval = errors.New("this is not the admin account")
		return errorval
	}
	return errorval
}
