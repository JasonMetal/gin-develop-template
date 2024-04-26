package main

import (
	"develop-template/constant"
	"github.com/JasonMetal/submodule-support-go.git/bootstrap"
	grpcHealth "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	bootstrap.SetProjectName(constant.ProjectName)
	bootstrap.Init()
	s := bootstrap.NewGrpcServer()
	// 健康检测
	grpc_health_v1.RegisterHealthServer(s, grpcHealth.NewServer())

	// 业务服务
	bootstrap.RunServer(s, constant.GrpcServiceHostPort)
}
