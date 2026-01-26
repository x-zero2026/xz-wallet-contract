# 快速启动指南

## 🚀 5 分钟快速体验

### 步骤 1: 安装依赖

```bash
npm install
```

### 步骤 2: 启动开发服务器

```bash
npm run dev
```

### 步骤 3: 测试登录

打开浏览器访问：**http://localhost:5173/test-login.html**

点击"模拟登录并跳转"按钮，即可进入任务中心。

### 步骤 4: 开始使用

现在你可以：
- 浏览可申请的任务
- 发布新任务
- 查看任务详情
- 申请任务
- 管理任务状态

## 📝 测试账号

默认测试账号（admin）：
- **用户名**: admin
- **DID**: 0x3070deb1c17432b094d30509ccbfd598fb2793435efdca9273dfbc558bc040ca
- **以太坊地址**: 0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA
- **XZT 余额**: 10,000 XZT
- **信用分**: 5,000

## 🔗 与 DID 登录系统集成

### 在 DID 登录系统中添加跳转按钮

编辑 `did-login-ui/src/views/Dashboard.jsx`：

```javascript
const navigateToTaskCenter = () => {
  window.location.href = 'http://localhost:5173';
};

// 在 Dashboard 中添加按钮
<button onClick={navigateToTaskCenter}>
  进入任务中心
</button>
```

### 配置环境变量

**DID 登录系统** (`.env`):
```env
VITE_TASK_CENTER_URL=http://localhost:5173
```

**任务中心** (`.env`):
```env
VITE_API_BASE_URL=https://yms07x0sn0.execute-api.us-east-1.amazonaws.com/prod
VITE_DID_LOGIN_URL=http://localhost:3000
```

## 🎯 主要功能

### 作为执行者
1. 浏览"可申请任务"标签页
2. 点击任务卡片查看详情
3. 点击"申请任务"按钮投标
4. 在"我的任务"中查看已申请的任务

### 作为发布者
1. 切换到"我发布的"标签页
2. 点击"+ 发布新任务"按钮
3. 填写任务信息并发布
4. 查看投标列表并选择执行者
5. 审批执行者提交的成果

## 🐛 常见问题

### Q: 显示"请先登录"
**A**: 访问 http://localhost:5173/test-login.html 进行测试登录

### Q: API 请求失败
**A**: 检查 `.env` 中的 `VITE_API_BASE_URL` 是否正确

### Q: 无法加载余额
**A**: 确保使用的是有效的 JWT Token 和以太坊地址

## 📚 更多文档

- [README.md](./README.md) - 完整文档
- [INTEGRATION.md](./INTEGRATION.md) - 集成指南
- [API 文档](../lambda/README.md) - 后端 API 文档

## 🎉 开始使用

现在你已经准备好了！访问 http://localhost:5173 开始体验任务中心。
