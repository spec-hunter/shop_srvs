package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/ahlemarg/shop-srvs/src/user_srvs/global"
	handler "github.com/ahlemarg/shop-srvs/src/user_srvs/handlers"
	"github.com/ahlemarg/shop-srvs/src/user_srvs/initialize"
	"github.com/ahlemarg/shop-srvs/src/user_srvs/proto"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	IP := flag.String("ip", "192.168.3.129", "ip地址")
	PORT := flag.Int("port", 50052, "prot端口")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *PORT))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	// 使用GRPC预留的proto, 给第三方的配置中心进行健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	c := api.DefaultConfig()
	// 需要将默认的地址进行更换
	// 默认地址为: 127.0.0.1:8500
	c.Address = fmt.Sprintf("%s:%d", global.ServerInfo.ConsulInfo.Host, global.ServerInfo.ConsulInfo.Port)

	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}

	s := new(api.AgentServiceRegistration)
	s.Name = global.ServerInfo.Name
	s.ID = global.ServerInfo.Name
	s.Tags = []string{"srvs", "user"}
	s.Address = "192.168.3.129"
	s.Port = 50052

	a := &api.AgentServiceCheck{
		// HTTP模式要求:  访问该URL能获得返回200状态码, 否则无法通过健康检查
		GRPC:                           "192.168.3.129:50052",
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "30s",
	}

	s.Check = a

	err = client.Agent().ServiceRegister(s)
	if err != nil {
		panic(err)
	}

	err = server.Serve(listen)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
