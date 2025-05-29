#!/bin/bash

# 订单管理系统停止脚本

echo "🛑 停止订单管理系统..."

STOPPED_SERVICES=0

# 停止后端服务
if [ -f backend.pid ]; then
    BACKEND_PID=$(cat backend.pid)
    if kill -0 $BACKEND_PID 2>/dev/null; then
        echo "🔴 停止后端服务 (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
        
        # 等待进程结束
        for i in {1..5}; do
            if ! kill -0 $BACKEND_PID 2>/dev/null; then
                break
            fi
            sleep 1
        done
        
        # 如果进程仍在运行，强制杀死
        if kill -0 $BACKEND_PID 2>/dev/null; then
            echo "⚠️  强制停止后端服务..."
            kill -9 $BACKEND_PID 2>/dev/null
        fi
        
        rm backend.pid
        STOPPED_SERVICES=$((STOPPED_SERVICES + 1))
        echo "✅ 后端服务已停止"
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
        
        # 等待进程结束
        for i in {1..3}; do
            if ! kill -0 $WEB_PID 2>/dev/null; then
                break
            fi
            sleep 1
        done
        
        # 如果进程仍在运行，强制杀死
        if kill -0 $WEB_PID 2>/dev/null; then
            echo "⚠️  强制停止Web服务..."
            kill -9 $WEB_PID 2>/dev/null
        fi
        
        rm web.pid
        STOPPED_SERVICES=$((STOPPED_SERVICES + 1))
        echo "✅ Web服务已停止"
    else
        echo "⚠️  Web服务已停止"
        rm -f web.pid
    fi
else
    echo "⚠️  未找到Web服务PID文件"
fi

# 清理其他可能的进程
echo "🧹 清理相关进程..."
KILLED_PROCESSES=0

# 清理后端进程
if pkill -f "uniapp-order-backend" 2>/dev/null; then
    KILLED_PROCESSES=$((KILLED_PROCESSES + 1))
    echo "   - 清理了遗留的后端进程"
fi

# 清理Web服务进程
if pkill -f "python3 -m http.server 8080" 2>/dev/null; then
    KILLED_PROCESSES=$((KILLED_PROCESSES + 1))
    echo "   - 清理了遗留的Web服务进程"
fi

# 清理Go运行进程
if pkill -f "go run.*main.go" 2>/dev/null; then
    KILLED_PROCESSES=$((KILLED_PROCESSES + 1))
    echo "   - 清理了遗留的Go运行进程"
fi

if [ $KILLED_PROCESSES -eq 0 ]; then
    echo "   - 没有发现遗留进程"
fi

# 清理临时文件
echo "🧹 清理临时文件..."
rm -f backend.pid web.pid 2>/dev/null

echo ""
if [ $STOPPED_SERVICES -gt 0 ] || [ $KILLED_PROCESSES -gt 0 ]; then
    echo "✅ 系统已完全停止 (停止了 $STOPPED_SERVICES 个服务，清理了 $KILLED_PROCESSES 个进程)"
else
    echo "✅ 系统已停止 (没有运行中的服务)"
fi

echo ""
echo "💡 提示: 使用 ./start.sh 可以重新启动系统"