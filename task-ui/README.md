# XZ 任务中心前端

基于 React + Vite 的任务管理系统前端应用。

## 功能特性

### 作为执行者（Executor）
- ✅ 浏览可申请的任务（招标中的任务）
- ✅ 申请任务（投标）
- ✅ 查看我的任务（已申请 + 进行中）
- ✅ 提交设计方案
- ✅ 提交实现代码
- ✅ 提交最终成果

### 作为发布者（Creator）
- ✅ 发布新任务
- ✅ 查看我发布的任务
- ✅ 查看任务的投标列表
- ✅ 选择执行者
- ✅ 审批设计方案（通过/拒绝）
- ✅ 审批实现代码（通过/拒绝）
- ✅ 审批最终成果（通过/拒绝，完成支付）
- ✅ 取消任务

### 通用功能
- ✅ 查看任务详情
- ✅ 查看 XZT 余额
- ✅ 查看信用分

## 技术栈

- **框架**: React 18
- **构建工具**: Vite
- **HTTP 客户端**: Axios
- **路由**: React Router DOM
- **样式**: CSS Modules

## 项目结构

```
task-ui/
├── src/
│   ├── api/
│   │   └── index.js          # API 接口定义
│   ├── components/
│   │   ├── TaskList.jsx      # 任务列表组件
│   │   ├── TaskCard.jsx      # 任务卡片组件
│   │   ├── CreateTaskModal.jsx    # 创建任务弹窗
│   │   └── TaskDetailModal.jsx    # 任务详情弹窗
│   ├── utils/
│   │   └── auth.js           # 认证工具函数
│   ├── App.jsx               # 主应用组件
│   ├── App.css               # 全局样式
│   └── main.jsx              # 应用入口
├── .env                      # 环境变量
├── .env.example              # 环境变量示例
└── package.json
```

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 配置环境变量

复制 `.env.example` 到 `.env` 并配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件：

```env
VITE_API_BASE_URL=https://yms07x0sn0.execute-api.us-east-1.amazonaws.com/prod
VITE_DID_LOGIN_URL=http://localhost:3000
```

### 3. 启动开发服务器

```bash
npm run dev
```

访问 http://localhost:5173

### 4. 测试登录

**方式 1: 使用测试页面**

访问 http://localhost:5173/test-login.html，填写测试用户信息并点击"模拟登录"。

**方式 2: 从 DID 登录系统跳转**

1. 启动 DID 登录系统
2. 登录成功后会自动跳转到任务中心
3. 任务中心会读取 localStorage 中的认证信息

**方式 3: 手动设置（浏览器控制台）**

```javascript
localStorage.setItem('token', 'your-jwt-token');
localStorage.setItem('userInfo', JSON.stringify({
  did: '0x...',
  username: 'admin',
  eth_address: '0x...'
}));
location.reload();
```

### 5. 构建生产版本

```bash
npm run build
```

构建产物在 `dist/` 目录。

## 任务状态流转

```
pending (待发布)
  ↓
bidding (招标中) ← 执行者可以申请
  ↓
accepted (已接受) ← 发布者选择执行者
  ↓
design_submitted (设计已提交) ← 执行者提交设计
  ↓
design_approved (设计已批准) ← 发布者批准设计，支付 30%
  ↓
implementation_submitted (实现已提交) ← 执行者提交实现
  ↓
implementation_approved (实现已批准) ← 发布者批准实现，支付 50%
  ↓
final_submitted (最终成果已提交) ← 执行者提交最终成果
  ↓
completed (已完成) ← 发布者批准最终成果，支付 20%
```

## API 接口

### 钱包相关
- `GET /wallet/balance?address={address}` - 查询 XZT 余额

### 任务相关
- `POST /tasks` - 创建任务
- `GET /tasks?status={status}&creator_did={did}&executor_did={did}` - 列出任务
- `GET /tasks/{id}` - 获取任务详情
- `POST /tasks/{id}/bid` - 投标任务
- `POST /tasks/{id}/select-bidder` - 选择投标者
- `POST /tasks/{id}/approve` - 审批工作

## 认证

应用使用 JWT Token 进行认证。Token 和用户信息存储在 localStorage 中。

### 从 DID 登录系统集成

任务中心不包含注册/登录页面，需要从 DID 登录系统跳转进入。

**DID 登录系统需要做的：**

登录成功后，设置 localStorage 并跳转：

```javascript
// 登录成功后
localStorage.setItem('token', response.token);
localStorage.setItem('userInfo', JSON.stringify({
  did: response.did,
  username: response.username,
  eth_address: response.eth_address,
  email: response.email,
  project_id: response.project_id || '00000000-0000-0000-0000-000000000000'
}));

// 跳转到任务中心
window.location.href = 'http://localhost:5173'; // 或生产环境 URL
```

**任务中心会自动读取：**

```javascript
const token = localStorage.getItem('token');
const userInfo = JSON.parse(localStorage.getItem('userInfo'));
```

所有 API 请求会自动在 Header 中添加：
```
Authorization: Bearer {token}
```

详细集成说明请查看 [INTEGRATION.md](./INTEGRATION.md)

## 开发说明

### 添加新功能

1. 在 `src/api/index.js` 中添加 API 接口
2. 在 `src/components/` 中创建新组件
3. 在 `App.jsx` 中集成新功能

### 样式规范

- 使用 CSS 类名，避免内联样式
- 遵循 BEM 命名规范
- 响应式设计，支持移动端

### 状态管理

当前使用 React 内置的 useState 和 useEffect。如果应用复杂度增加，可以考虑：
- Context API
- Redux
- Zustand

## 部署

### 部署到 AWS Amplify

1. 创建 `amplify.yml`：

```yaml
version: 1
frontend:
  phases:
    preBuild:
      commands:
        - npm ci
    build:
      commands:
        - npm run build
  artifacts:
    baseDirectory: dist
    files:
      - '**/*'
  cache:
    paths:
      - node_modules/**/*
```

2. 在 Amplify Console 中连接 Git 仓库
3. 配置环境变量
4. 部署

### 部署到其他平台

- **Vercel**: `vercel --prod`
- **Netlify**: `netlify deploy --prod`
- **静态服务器**: 将 `dist/` 目录部署到任何静态服务器

## 故障排查

### Token 过期
- 检查 localStorage 中的 token
- 重新登录获取新 token

### API 请求失败
- 检查 `.env` 中的 API_BASE_URL
- 检查网络连接
- 查看浏览器控制台错误信息

### 样式问题
- 清除浏览器缓存
- 检查 CSS 文件是否正确导入

## License

MIT
