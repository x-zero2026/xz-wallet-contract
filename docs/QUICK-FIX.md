# 快速修复 - Token 处理

## 问题
从 DID 登录跳转时，URL 中的 token 无法正确处理，导致页面报错。

## 已修复 ✅

### 1. Token 处理逻辑
- ✅ 从 URL 参数读取 token
- ✅ 保存到 localStorage
- ✅ 解析 JWT 获取用户信息
- ✅ 清除 URL 中的 token
- ✅ 自动刷新页面

### 2. ETH 地址映射
由于 DID Login API 的 `/api/profile` 端点也有认证问题，暂时使用硬编码映射：

```javascript
const knownUsers = {
  '0x3070deb1c17432b094d30509ccbfd598fb2793435efdca9273dfbc558bc040ca': 
    '0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA'  // admin
};
```

## 使用方法

### 从 DID 登录系统跳转
1. 访问 http://localhost:3000
2. 登录 admin 账号
3. 点击 "人才市场" 💰 应用
4. 自动跳转到 http://localhost:5173/?token=xxx
5. ✅ 自动登录并显示任务中心

### 直接访问测试
访问：
```
http://localhost:5173/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkaWQiOiIweDMwNzBkZWIxYzE3NDMyYjA5NGQzMDUwOWNjYmZkNTk4ZmIyNzkzNDM1ZWZkY2E5MjczZGZiYzU1OGJjMDQwY2EiLCJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzcwMDMzMjA4LCJpYXQiOjE3Njk0Mjg0MDh9.WegEA-vWKBUnj9DPR8hwAJez76sRZuAU0suA49Q9384
```

## 后续改进

### 需要修复 DID Login API Gateway
DID Login 的 API 也需要添加 `Auth: DefaultAuthorizer: NONE`，然后可以：

1. 移除硬编码的用户映射
2. 从 `/api/profile` 端点获取完整用户信息
3. 支持所有用户（不仅仅是 admin）

### 修复步骤
在 `did-login-lambda/template.yaml` 中添加：

```yaml
ApiGateway:
  Type: AWS::Serverless::Api
  Properties:
    StageName: prod
    Auth:
      DefaultAuthorizer: NONE  # 添加这一行
    Cors:
      # ... 其他配置
```

然后重新部署：
```bash
cd did-login-lambda
sam build && sam deploy
```

## 当前状态

- ✅ Admin 用户可以正常登录
- ✅ 可以查看任务列表
- ✅ 可以查看余额
- ⚠️ 其他用户需要在代码中添加映射

---

**更新时间**: 2026-01-26  
**状态**: ✅ 临时修复完成
