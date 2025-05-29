# 订单管理系统 - UniApp前端

## 项目说明

这是订单管理系统的UniApp前端部分，支持编译到多个平台：
- H5网页
- 微信小程序
- 支付宝小程序
- App（Android/iOS）

## 目录结构

```
frontend/
├── pages/              # 页面
│   ├── dashboard/      # 驾驶仓
│   ├── orders/         # 订单管理
│   ├── products/       # 商品管理
│   ├── customers/      # 客户管理
│   └── calendar/       # 日历视图
├── components/         # 组件
├── utils/              # 工具函数
│   └── api.js          # API接口
├── static/             # 静态资源
├── App.vue            # 应用入口
├── main.js            # 主程序
├── pages.json         # 页面配置
└── manifest.json      # 应用配置
```

## 开发环境

### 使用HBuilderX（推荐）

1. 下载并安装 [HBuilderX](https://www.dcloud.io/hbuilderx.html)
2. 打开HBuilderX，选择"文件" -> "打开目录"
3. 选择当前frontend目录
4. 在工具栏选择运行平台（浏览器、微信开发者工具等）
5. 点击"运行"按钮

### 使用CLI工具

```bash
# 安装uni-app CLI
npm install -g @vue/cli @vue/cli-init
vue init dcloudio/uni-preset-vue my-project

# 或使用新版CLI
npm install -g @dcloudio/uvm
uvm install latest
npm install -g @dcloudio/cli
```

## 配置说明

### API地址配置

编辑 `utils/api.js` 文件中的 `BASE_URL`：

```javascript
const BASE_URL = 'http://localhost:12000/api'  // 开发环境
// const BASE_URL = 'https://your-domain.com/api'  // 生产环境
```

### 应用配置

编辑 `manifest.json` 文件配置应用信息：
- 应用名称
- 应用ID
- 版本信息
- 平台特定配置

## 页面功能

### 1. 驾驶仓 (dashboard)
- 显示总体数据统计
- 今日/本月订单概况
- 客户和商品统计

### 2. 订单管理 (orders)
- 订单列表查看
- 新增订单
- 编辑订单
- 订单状态管理

### 3. 商品管理 (products)
- 商品列表
- 添加商品
- 编辑商品信息
- 商品分类管理

### 4. 客户管理 (customers)
- 客户列表
- 客户信息管理
- 多地址管理

### 5. 日历视图 (calendar)
- 月历显示
- 农历信息
- 订单统计
- 节假日显示

## 编译发布

### H5网页
```bash
# 在HBuilderX中选择"发行" -> "网站-H5手机版"
```

### 微信小程序
```bash
# 在HBuilderX中选择"发行" -> "小程序-微信"
# 需要先配置微信小程序AppID
```

### App
```bash
# 在HBuilderX中选择"发行" -> "原生App-云打包"
# 需要配置证书和打包参数
```

## 注意事项

1. 确保后端API服务正常运行
2. 跨域问题：开发时可能需要配置代理
3. 小程序发布需要相应平台的开发者账号
4. App打包需要配置相应的证书和签名

## 技术栈

- **框架**: UniApp
- **语言**: Vue.js + JavaScript
- **UI**: 原生组件 + 自定义样式
- **网络**: uni.request API
- **存储**: uni.storage API