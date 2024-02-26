package initialize

import (
	"github.com/gin-gonic/gin"
	"go-apis/user/middlewares"
	userRouter "go-apis/user/router"
)

// Routers 初始化路由
func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/user/v1")
	userRouter.InitUserRouter(ApiGroup)
	userRouter.InitBaseRouter(ApiGroup)
	return Router
}
