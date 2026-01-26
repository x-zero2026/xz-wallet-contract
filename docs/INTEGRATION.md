# 任务中心集成指南

## 从 DID 登录系统集成

任务中心需要从 DID 登录系统接收用户认证信息。

### 1. 登录后跳转

在 DID 登录系统登录成功后，跳转到任务中心并传递认证信息：

```javascript
// 登录成功后
const loginResponse = {
  token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  did: "0x3070deb1c17432b094d30509ccbfd598fb2793435efdca9273dfbc558bc040ca",
  username: "admin",
  eth_address: "0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA",
  email: "x-zero2026@outlook.com"
};

// 存储到 localStorage
localStorage.setItem('token', loginResponse.token);
localStorage.setItem('userInfo', JSON.stringify({
  did: loginResponse.did,
  username: loginResponse.username,
  eth_address: loginResponse.eth_address,
  email: loginResponse.email,
  project_id: loginResponse.project_id || '00000000-0000-0000-0000-000000000000'
}));

// 跳转到任务中心
window.location.href = 'http://localhost:5173'; // 或生产环境 URL
```

### 2. 所需的用户信息

任务中心需要以下用户信息：

```typescript
interface UserInfo {
  did: string;              // 用户 DID（必需）
  username: string;         // 用户名（必需）
  eth_address: string;      // 以太坊地址（必需）
  email?: string;           // 邮箱（可选）
  project_id?: string;      // 项目 ID（可选，默认使用全局项目）
}
```

### 3. JWT Token 格式

Token 应该包含以下 claims：

```json
{
  "did": "0x...",
  "username": "admin",
  "exp": 1769513586,
  "iat": 1769427186
}
```

### 4. 示例：完整的登录流程

#### DID 登录页面（did-login-ui）

```javascript
// 在 Dashboard.jsx 或登录成功后
const handleLoginSuccess = (response) => {
  // 存储认证信息
  localStorage.setItem('token', response.token);
  localStorage.setItem('userInfo', JSON.stringify({
    did: response.did,
    username: response.username,
    eth_address: response.eth_address,
    email: response.email
  }));
  
  // 跳转到任务中心
  window.location.href = process.env.VITE_TASK_CENTER_URL || 'http://localhost:5173';
};
```

#### 任务中心（task-ui）

```javascript
// App.jsx 会自动读取 localStorage 中的认证信息
useEffect(() => {
  const user = getUserInfo(); // 从 localStorage 读取
  if (!user) {
    // 显示未登录提示
    return;
  }
  setUserInfo(user);
  loadBalance(user.eth_address);
}, []);
```

### 5. 退出登录

任务中心的退出按钮会清除本地存储并刷新页面：

```javascript
const handleLogout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('userInfo');
  window.location.reload(); // 或跳转回登录页
};
```

如果需要跳转回 DID 登录页面：

```javascript
const handleLogout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('userInfo');
  window.location.href = process.env.VITE_DID_LOGIN_URL || 'http://localhost:5174';
};
```

### 6. 环境变量配置

#### DID 登录系统 (.env)
```env
VITE_TASK_CENTER_URL=http://localhost:5173
```

#### 任务中心 (.env)
```env
VITE_API_BASE_URL=https://yms07x0sn0.execute-api.us-east-1.amazonaws.com/prod
VITE_DID_LOGIN_URL=http://localhost:5174
```

### 7. 测试集成

#### 方式 1: 手动设置（用于测试）

在浏览器控制台执行：

```javascript
// 设置测试用户
localStorage.setItem('token', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkaWQiOiIweDMwNzBkZWIxYzE3NDMyYjA5NGQzMDUwOWNjYmZkNTk4ZmIyNzkzNDM1ZWZkY2E5MjczZGZiYzU1OGJjMDQwY2EiLCJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzY5NTEzNTg2LCJpYXQiOjE3Njk0MjcxODZ9.xFST95uXdhMYbCy7p_B9BVR1FRqPzeDEDHPxNY8Rv4Q');

localStorage.setItem('userInfo', JSON.stringify({
  did: '0x3070deb1c17432b094d30509ccbfd598fb2793435efdca9273dfbc558bc040ca',
  username: 'admin',
  eth_address: '0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA',
  email: 'x-zero2026@outlook.com'
}));

// 刷新页面
location.reload();
```

#### 方式 2: 从 DID 登录系统跳转

1. 启动 DID 登录系统：`cd did-login-ui && npm run dev`
2. 启动任务中心：`cd task-ui && npm run dev`
3. 在 DID 登录系统登录
4. 登录成功后自动跳转到任务中心

### 8. 生产环境部署

#### 部署到 AWS Amplify

两个应用都部署到 Amplify 后：

1. **DID 登录系统**
   - URL: `https://did-login.example.com`
   - 环境变量: `VITE_TASK_CENTER_URL=https://task-center.example.com`

2. **任务中心**
   - URL: `https://task-center.example.com`
   - 环境变量: 
     - `VITE_API_BASE_URL=https://yms07x0sn0.execute-api.us-east-1.amazonaws.com/prod`
     - `VITE_DID_LOGIN_URL=https://did-login.example.com`

### 9. 安全注意事项

1. **Token 过期处理**
   - 任务中心会自动检测 401 错误
   - 可以配置自动跳转回登录页

2. **HTTPS**
   - 生产环境必须使用 HTTPS
   - 确保 Token 传输安全

3. **CORS**
   - API Gateway 已配置 CORS
   - 允许所有来源（生产环境建议限制）

4. **Token 刷新**
   - 当前实现不支持自动刷新
   - Token 过期后需要重新登录

### 10. 故障排查

#### 问题：任务中心显示"请先登录"

**原因**：localStorage 中没有认证信息

**解决**：
1. 检查是否从 DID 登录系统正确跳转
2. 检查浏览器控制台是否有错误
3. 手动设置测试数据（见上文）

#### 问题：API 请求返回 401

**原因**：Token 无效或过期

**解决**：
1. 重新登录获取新 Token
2. 检查 Token 格式是否正确
3. 验证 JWT_SECRET 配置

#### 问题：无法加载余额

**原因**：eth_address 不正确或 API 错误

**解决**：
1. 检查 userInfo 中的 eth_address
2. 检查 API_BASE_URL 配置
3. 查看浏览器网络请求

## 完整示例代码

### DID 登录系统集成代码

在 `did-login-ui/src/views/Dashboard.jsx` 中添加：

```javascript
const navigateToTaskCenter = () => {
  const taskCenterUrl = import.meta.env.VITE_TASK_CENTER_URL || 'http://localhost:5173';
  window.location.href = taskCenterUrl;
};

// 在 Dashboard 组件中添加按钮
<button onClick={navigateToTaskCenter}>
  进入任务中心
</button>
```

### 任务中心返回登录页

在 `task-ui/src/App.jsx` 中修改退出逻辑：

```javascript
const handleLogout = () => {
  logout();
  const loginUrl = import.meta.env.VITE_DID_LOGIN_URL || 'http://localhost:5174';
  window.location.href = loginUrl;
};
```
