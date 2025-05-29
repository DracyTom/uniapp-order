#!/bin/bash

# 订单管理系统重启脚本

echo "🔄 重启订单管理系统..."

# 停止现有服务
echo "🛑 停止现有服务..."
./stop.sh

# 等待一下确保进程完全停止
sleep 2

# 启动服务
echo ""
echo "🚀 启动服务..."
./start.sh