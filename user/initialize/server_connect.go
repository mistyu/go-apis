package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go-apis/user/global"
	"go-apis/user/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitServerConnect() {
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[InitServerConnect] 连接【用户服务失败】")
	}

	// 后续的服务下线了; 改端口了; 改ip了 负载均衡来做
	// 事先已经创立好了连接，这样后续就不用进行多次tcp的三次握手
	// 一个连接多个groutine公用，性能 - 连接池
	userServerClient := proto.NewUserClient(userConn)
	global.UserServerClient = userServerClient

}

//func InitServerConnect1() {
//	// 从注册中心获取到用户信息
//	cfg := api.DefaultConfig()
//	consulInfo := global.ServerConfig.ConsulInfo
//	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)
//	userServerHost := ""
//	userServerPort := 0
//	client, err := api.NewClient(cfg)
//	if err != nil {
//		panic(err)
//	}
//	data, err := client.Agent().ServicesWithFilter("")
//	if err != nil {
//		panic(err)
//	}
//	for _, value := range data {
//		userServerHost = value.Address
//		userServerPort = value.Port
//		break
//	}
//	if userServerHost == "" {
//		zap.S().Fatal("[InitServerConnect] 连接【用户服务失败】")
//		return
//	}
//	// 拨号连接用户grpc服务器
//	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userServerHost, userServerPort))
//	if err != nil {
//		zap.S().Errorw("[GetUserList] 连接【用户服务失败】", "msg", err.Error())
//	}
//
//	// 后续的服务下线了; 改端口了; 改ip了 负载均衡来做
//	// 事先已经创立好了连接，这样后续就不用进行多次tcp的三次握手
//	// 一个连接多个groutine公用，性能 - 连接池
//	userServerClient := proto.NewUserClient(userConn)
//	global.UserServerClient = userServerClient
//}
