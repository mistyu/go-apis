package router

import (
	"github.com/gin-gonic/gin"
	"go-apis/user/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	{
		userRouter.GET("list", api.GetUserList)
		userRouter.POST("password_login", api.PasswordLogin)
	}
}
