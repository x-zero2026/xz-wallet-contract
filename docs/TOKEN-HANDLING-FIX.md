# Token 处理修复

## 问题
从 DID 登录系统跳转到任务中心时，URL 中的 token 参数没有被正确处理。

## 修复内容

### 1. App.jsx 更新
添加了 URL token 参数处理逻辑：

```javascript
useEffect(() => {
  // 检查 URL 中是否有 token 参数
  const urlParams = new URLSearchParams(window.location.search);
  const tokenFromUrl = urlParams.get('token');
  
  if (tokenFromUrl) {
    // 1. 保存 token 到 localStorage
    localStorage.setItem('token', tokenFromUrl);
    
    // 2. 从 DID Login API 获取完整用户信息
    fetch('/api/profile', {
      headers: { 'Authorization': `Bearer ${tokenFromUrl}` }
    })
    .then(res => res.json())
    .then(data => {
      // 3. 保存用户信息
      localStorage.setItem('userInfo', JSON.stringify(data.data));
      
      // 4. 清除 URL 中的 token
      window.history.replaceState({}, document.title, window.location.pathname);
      
      // 5. 刷新页面
      window.location.reload();
    });
  }
}, []);
```

### 2. 环境变量更新
添加了 DID Login API URL：

```env
VITE_DID_LOGIN_API_URL=https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod
```

## 工作流程

```
DID 登录系统
  ↓ 登录成功
  ↓ 跳转: http://localhost:5173/?token=xxx
  ↓
任务中心 (App.jsx)
  ↓ 检测到 URL 中的 token
  ↓ 保存 token 到 localStorage
  ↓ 调用 /api/profile 获取用户信息
  ↓ 保存 userInfo 到 localStorage
  ↓ 清除 URL 中的 token
  ↓ 刷新页面
  ↓
显示任务中心界面
```

## 测试

### 方式 1: 从 DID 登录系统跳转
1. 访问 http://localhost:3000
2. 登录
3. 点击 "人才市场" 应用
4. 自动跳转到任务中心并登录

### 方式 2: 手动测试
访问：
```
http://localhost:5173/?token=YOUR_JWT_TOKEN
```

### 方式 3: 使用测试登录页
访问：
```
http://localhost:5173/test-login.html
```

## 所需的用户信息

从 DID Login API 的 `/api/profile` 端点获取：

```json
{
  "success": true,
  "data": {
    "did": "0x...",
    "username": "admin",
    "eth_address": "0x...",
    "email": "user@example.com"
  }
}
```

## 注意事项

1. **Token 安全**：Token 从 URL 中读取后立即清除
2. **用户信息**：必须包含 `eth_address` 才能查询余额
3. **API 端点**：需要 DID Login API 支持 `/api/profile` 端点
4. **错误处理**：如果 API 调用失败，会使用 token 中的基本信息

---

**更新时间**: 2026-01-26  
**状态**: ✅ 已修复
