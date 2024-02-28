package router

import (
	"github.com/gin-gonic/gin"
	"go-apis/user/api"
	"go-apis/user/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	{
		userRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		userRouter.POST("password_login", api.PasswordLogin)
		userRouter.POST("register", api.Register)
	}
}
