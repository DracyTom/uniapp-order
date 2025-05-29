#!/bin/bash

# 订单管理系统启动脚本

echo "🚀 启动订单管理系统..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go环境未安装，请先安装Go"
    exit 1
fi

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
sleep 3

# 检查后端是否启动成功
if curl -s http://localhost:12000/api/products > /dev/null; then
    echo "✅ 后端服务启动成功"
else
    echo "❌ 后端服务启动失败"
    kill $BACKEND_PID 2>/dev/null
    exit 1
fi

# 返回项目根目录
cd ..

# 启动测试页面服务器
echo "🌐 启动测试页面服务器 (端口: 8080)..."
python3 -m http.server 8080 > web_server.log 2>&1 &
WEB_PID=$!

echo ""
echo "🎉 系统启动完成！"
echo ""
echo "📊 访问地址："
echo "   - API测试页面: http://localhost:8080/test.html"
echo "   - 后端API: http://localhost:12000/api"
echo ""
echo "📝 日志文件："
echo "   - 后端日志: backend/server.log"
echo "   - Web服务日志: web_server.log"
echo ""
echo "🛑 停止服务："
echo "   - 后端PID: $BACKEND_PID"
echo "   - Web服务PID: $WEB_PID"
echo "   - 停止命令: kill $BACKEND_PID $WEB_PID"
echo ""

# 保存PID到文件
echo $BACKEND_PID > backend.pid
echo $WEB_PID > web.pid

echo "💡 提示: 使用 ./stop.sh 可以停止所有服务"