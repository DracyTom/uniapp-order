package main

import (
	"log"

	"uniapp-order-backend/config"
	"uniapp-order-backend/models"
	"uniapp-order-backend/routes"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	models.InitDB()

	// 设置路由
	r := routes.SetupRoutes()

	// 启动服务器
	addr := config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
	log.Printf("服务器启动在: %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}