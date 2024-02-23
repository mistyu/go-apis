package initialize

import (
	"github.com/gin-gonic/gin"
	userRouter "go-apis/user/router"
)

// Routers 初始化路由
func Routers() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("/user/v1")
	userRouter.InitUserRouter(ApiGroup)
	return Router
}
