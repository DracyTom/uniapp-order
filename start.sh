#!/bin/bash

# 订单管理系统启动脚本

echo "🚀 启动订单管理系统..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go环境未安装，请先安装Go"
    exit 1
fi

# 检查端口是否被占用
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo "❌ 端口 $port 已被占用"
        echo "请先停止占用端口的进程或使用 ./stop.sh 停止现有服务"
        exit 1
    fi
}

echo "🔍 检查端口占用情况..."
check_port 12000
check_port 8080

# 进入后端目录
cd backend

# 检查依赖
echo "📦 检查Go依赖..."
go mod tidy

# 编译后端
echo "🔨 编译后端程序..."
go build -o uniapp-order-backend main.go

if [ $? -ne 0 ]; then
    echo "❌ 后端编译失败"
    exit 1
fi

# 启动后端服务
echo "🌐 启动后端服务 (端口: 12000)..."
./uniapp-order-backend > server.log 2>&1 &
BACKEND_PID=$!

# 等待后端启动
echo "⏳ 等待后端服务启动..."
sleep 5

# 检查后端是否启动成功
echo "🔍 检查后端服务状态..."
for i in {1..10}; do
    if curl -s http://localhost:12000/api/dashboard > /dev/null 2>&1; then
        echo "✅ 后端服务启动成功"
        break
    elif [ $i -eq 10 ]; then
        echo "❌ 后端服务启动失败，请检查日志: backend/server.log"
        kill $BACKEND_PID 2>/dev/null
        exit 1
    else
        echo "⏳ 等待后端服务响应... ($i/10)"
        sleep 2
    fi
done

# 返回项目根目录
cd ..

# 启动测试页面服务器
echo "🌐 启动测试页面服务器 (端口: 8080)..."
if command -v python3 &> /dev/null; then
    python3 -m http.server 8080 > web_server.log 2>&1 &
    WEB_PID=$!
    
    # 等待Web服务器启动
    sleep 2
    
    # 检查Web服务器是否启动成功
    if curl -s http://localhost:8080/ > /dev/null 2>&1; then
        echo "✅ Web服务器启动成功"
    else
        echo "❌ Web服务器启动失败"
        kill $BACKEND_PID $WEB_PID 2>/dev/null
        exit 1
    fi
else
    echo "⚠️  Python3未安装，跳过Web服务器启动"
    echo "   您可以手动启动: python3 -m http.server 8080"
    WEB_PID=""
fi

echo ""
echo "🎉 系统启动完成！"
echo ""
echo "📊 访问地址："
echo "   - 后端API: http://localhost:12000/api/dashboard"
if [ -n "$WEB_PID" ]; then
    echo "   - API测试页面: http://localhost:8080/test.html"
    echo "   - 前端演示页面: http://localhost:8080/frontend-test.html"
fi
echo ""
echo "📝 日志文件："
echo "   - 后端日志: backend/server.log"
if [ -n "$WEB_PID" ]; then
    echo "   - Web服务日志: web_server.log"
fi
echo ""
echo "🛑 停止服务："
echo "   - 后端PID: $BACKEND_PID"
if [ -n "$WEB_PID" ]; then
    echo "   - Web服务PID: $WEB_PID"
    echo "   - 停止命令: kill $BACKEND_PID $WEB_PID"
else
    echo "   - 停止命令: kill $BACKEND_PID"
fi
echo ""

# 保存PID到文件
echo $BACKEND_PID > backend.pid
if [ -n "$WEB_PID" ]; then
    echo $WEB_PID > web.pid
fi

echo "💡 提示: 使用 ./stop.sh 可以停止所有服务"
echo ""
echo "🔧 开发提示:"
echo "   - 后端代码修改后需要重新编译: cd backend && go build -o uniapp-order-backend main.go"
echo "   - 前端开发请使用 HBuilderX 或 uni-app CLI 工具"
echo "   - 生产环境建议使用 PostgreSQL 数据库"