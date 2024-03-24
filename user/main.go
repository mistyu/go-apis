package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go-apis/user/global"
	"go-apis/user/initialize"
	"go-apis/user/utils"
	myValidator "go-apis/user/validator"
	"go.uber.org/zap"
)

func main() {
	// 初始化Routers
	Router := initialize.Routers()
	// 初始化logger
	initialize.InitLogger()
	// 初始化配置文件
	initialize.InitConfig()
	// 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	// 初始化 server 的连接
	initialize.InitServerConnect()
	port := global.ServerConfig.Port
	viper.AutomaticEnv()
	// 如果是本地开发环境就固定
	//debug := viper.GetBool("GO_DEBUG")
	//if debug {
	//	port, err := utils.GetFreePort()
	//	if err == nil {
	//		global.ServerConfig.Port = port
	//	}
	//}

	port, err := utils.GetFreePort()
	if err == nil {
		global.ServerConfig.Port = port
	}

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("mobile", myValidator.ValidateMobile)
		if err != nil {
			panic(err)
		}
		err = v.RegisterTranslation("required", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 手机号码格式不对!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
		if err != nil {
			panic(err)
		}
	}

	// zap.S()可以获取一个全局的sugar
	zap.S().Infof("启动服务器，端口：%d", port)
	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
