#!/bin/bash

# 订单管理系统停止脚本

echo "🛑 停止订单管理系统..."

# 停止后端服务
if [ -f backend.pid ]; then
    BACKEND_PID=$(cat backend.pid)
    if kill -0 $BACKEND_PID 2>/dev/null; then
        echo "🔴 停止后端服务 (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
        rm backend.pid
    else
        echo "⚠️  后端服务已停止"
        rm -f backend.pid
    fi
else
    echo "⚠️  未找到后端服务PID文件"
fi

# 停止Web服务
if [ -f web.pid ]; then
    WEB_PID=$(cat web.pid)
    if kill -0 $WEB_PID 2>/dev/null; then
        echo "🔴 停止Web服务 (PID: $WEB_PID)..."
        kill $WEB_PID
        rm web.pid
    else
        echo "⚠️  Web服务已停止"
        rm -f web.pid
    fi
else
    echo "⚠️  未找到Web服务PID文件"
fi

# 清理其他可能的进程
echo "🧹 清理相关进程..."
pkill -f "uniapp-order-backend" 2>/dev/null
pkill -f "python3 -m http.server 8080" 2>/dev/null

echo "✅ 系统已停止"