package initialize

import "go.uber.org/zap"

// InitLogger 初始化logger
func InitLogger() {
	// 初始化日志
	logger, _ := zap.NewDevelopment()
	// 取代 zap 的logger
	zap.ReplaceGlobals(logger)
}
