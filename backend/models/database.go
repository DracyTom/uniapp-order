package models

import (
	"fmt"
	"log"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"uniapp-order-backend/config"
)

var DB *xorm.Engine

func InitDB() {
	var err error
	cfg := config.AppConfig.Database

	var dataSource string
	switch cfg.Driver {
	case "sqlite3":
		dataSource = cfg.Database
	case "postgres":
		dataSource = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)
	default:
		log.Fatalf("不支持的数据库驱动: %s", cfg.Driver)
	}

	DB, err = xorm.NewEngine(cfg.Driver, dataSource)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 设置连接池
	DB.SetMaxIdleConns(10)
	DB.SetMaxOpenConns(100)

	// 同步数据库结构
	err = DB.Sync2(
		new(Product),
		new(Customer),
		new(CustomerAddress),
		new(Order),
		new(OrderItem),
		new(Holiday),
	)
	if err != nil {
		log.Fatalf("数据库同步失败: %v", err)
	}

	log.Println("数据库初始化成功")
}