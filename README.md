# 订单管理系统

一个基于 Go + UniApp 的完整订单管理系统，支持商品管理、客户管理、订单录入、日历视图和数据统计。

## 🚀 项目状态

**✅ 已完成功能：**
- 完整的后端API系统（Go + Gin + XORM）
- 前端页面和组件（UniApp + Vue.js）
- 数据库设计和自动迁移
- API接口测试页面
- 项目文档和部署指南

**🔧 当前版本：** v1.0.0  
**📅 最后更新：** 2025-05-29

## 功能特性

### 1. 商品管理
- ✅ 商品列表展示（支持分页、搜索、筛选）
- ✅ 商品信息管理（名称、类别、价格、上架状态、图片）
- ✅ 自定义属性支持
- ✅ 商品分类管理

### 2. 客户管理
- ✅ 客户信息管理（姓名、联系方式）
- ✅ 多地址管理（支持默认地址设置）
- ✅ 客户搜索功能

### 3. 订单管理
- ✅ 订单录入界面
- ✅ 送货时间、地址、收货人管理
- ✅ 商品选择和数量管理
- ✅ 自动计算订单金额
- ✅ 支付状态管理（已付款/未付款）
- ✅ 特殊要求和备注
- ✅ 订单状态管理（待确认/已确认/已送达/已取消）

### 4. 日历视图
- ✅ 年/月/周三种视图模式
- ✅ 公历显示
- ✅ 农历显示（简化版）
- ✅ 节假日显示
- ✅ 每日订单量统计
- ✅ 订单金额统计

### 5. 驾驶仓页面
- ✅ 今日/本月/本年订单统计
- ✅ 订单状态分布
- ✅ 近7天订单趋势
- ✅ 热销商品排行
- ✅ 客户统计信息

### 6. 数据存储
- ✅ 支持 SQLite（测试环境）
- ✅ 支持 PostgreSQL（生产环境）
- ✅ 统一的数据库配置管理

## 技术栈

### 后端
- **Go 1.21+** - 编程语言
- **Gin** - Web框架
- **XORM** - ORM框架
- **SQLite/PostgreSQL** - 数据库
- **Viper** - 配置管理

### 前端
- **UniApp** - 跨平台开发框架
- **Vue.js** - 前端框架
- **uni-ui** - UI组件库

## 项目结构

```
uniapp-order/
├── backend/                 # Go后端
│   ├── config/             # 配置管理
│   │   ├── config.go
│   │   └── config.yaml
│   ├── models/             # 数据模型
│   │   ├── models.go
│   │   └── database.go
│   ├── controllers/        # 控制器
│   │   ├── product.go
│   │   ├── customer.go
│   │   ├── order.go
│   │   └── dashboard.go
│   ├── routes/             # 路由配置
│   │   └── routes.go
│   ├── utils/              # 工具函数
│   │   ├── response.go
│   │   └── lunar.go
│   ├── go.mod
│   └── main.go
├── frontend/               # UniApp前端
│   ├── pages/              # 页面
│   │   ├── dashboard/      # 驾驶仓
│   │   ├── products/       # 商品管理
│   │   ├── customers/      # 客户管理
│   │   ├── orders/         # 订单管理
│   │   └── calendar/       # 日历视图
│   ├── components/         # 组件
│   ├── utils/              # 工具函数
│   │   ├── api.js          # API接口
│   │   └── common.js       # 通用函数
│   ├── static/             # 静态资源
│   ├── App.vue
│   ├── main.js
│   ├── pages.json
│   └── manifest.json
└── README.md
```

## 快速开始

### 🚀 一键启动（推荐）

```bash
# 启动所有服务
./start.sh

# 停止所有服务
./stop.sh
```

启动后访问：
- **API测试页面**: http://localhost:8080/test.html
- **后端API**: http://localhost:12000/api

### 🧪 测试验证

系统提供了完整的API测试页面，可以验证所有功能：

1. 访问 http://localhost:8080/test.html
2. 点击各个功能按钮测试API
3. 查看返回的JSON数据

测试功能包括：
- ✅ 驾驶仓数据统计
- ✅ 商品管理（创建、查询）
- ✅ 客户管理（创建、查询）
- ✅ 订单管理（创建、查询）
- ✅ 日历数据（含农历、订单统计）

### 环境要求
- Go 1.21+
- Node.js 16+
- HBuilderX（用于UniApp开发）

### 手动启动

#### 后端启动

1. 进入后端目录
```bash
cd backend
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置数据库（可选）
编辑 `config/config.yaml` 文件，默认使用SQLite

4. 启动服务
```bash
go run main.go
```

后端服务将在 `http://localhost:12000` 启动

#### 前端启动

1. 使用HBuilderX打开 `frontend` 目录

2. 配置API地址
编辑 `utils/api.js` 中的 `BASE_URL`

3. 运行到浏览器或模拟器

## API接口

### 商品管理
- `GET /api/products` - 获取商品列表
- `GET /api/products/:id` - 获取单个商品
- `POST /api/products` - 创建商品
- `PUT /api/products/:id` - 更新商品
- `DELETE /api/products/:id` - 删除商品
- `GET /api/products/categories` - 获取商品分类

### 客户管理
- `GET /api/customers` - 获取客户列表
- `GET /api/customers/:id` - 获取单个客户
- `POST /api/customers` - 创建客户
- `PUT /api/customers/:id` - 更新客户
- `DELETE /api/customers/:id` - 删除客户
- `GET /api/customers/:id/addresses` - 获取客户地址
- `POST /api/customers/:id/addresses` - 添加客户地址

### 订单管理
- `GET /api/orders` - 获取订单列表
- `GET /api/orders/:id` - 获取单个订单
- `POST /api/orders` - 创建订单
- `PUT /api/orders/:id` - 更新订单
- `DELETE /api/orders/:id` - 删除订单
- `GET /api/orders/calendar` - 获取日历数据

### 驾驶仓
- `GET /api/dashboard` - 获取驾驶仓数据

## 数据库配置

### SQLite（默认）
```yaml
database:
  driver: "sqlite3"
  database: "./data.db"
```

### PostgreSQL
```yaml
database:
  driver: "postgres"
  host: "localhost"
  port: "5432"
  username: "your_username"
  password: "your_password"
  database: "order_management"
  sslmode: "disable"
```

## 部署说明

### 后端部署
1. 编译二进制文件
```bash
go build -o order-backend main.go
```

2. 配置生产环境数据库

3. 启动服务
```bash
./order-backend
```

### 前端部署
1. 使用HBuilderX发行到各个平台
2. 或者编译为H5版本部署到Web服务器

## 开发计划

- [ ] 图片上传功能
- [ ] 订单打印功能
- [ ] 数据导出功能
- [ ] 用户权限管理
- [ ] 消息推送
- [ ] 数据备份恢复
- [ ] 更完整的农历和节假日支持

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

如有问题或建议，请提交 Issue 或联系开发者。