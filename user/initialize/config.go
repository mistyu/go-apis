package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go-apis/user/global"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("GO_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user/%s.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user/%s-debug.yaml", configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Info("配置信息: &v", global.ServerConfig)
	// 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
	})
}
