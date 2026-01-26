# 自动 Approve 实现方案

## 概述
当用户发布任务时，如果尚未 approve escrow 合约，系统会自动调用 did-login-lambda 的 `approve-escrow` API 来完成 approval，无需用户手动操作。

## 架构

```
用户发布任务
    ↓
xz-wallet create-task Lambda
    ↓
检查 allowance
    ↓
如果不足 → 调用 did-login approve-escrow API
    ↓
approve-escrow Lambda:
  1. 从 Vault 获取助记词
  2. 派生以太坊私钥
  3. 执行 approve 交易
  4. 返回成功（不返回私钥）
    ↓
create-task 继续创建任务
```

## 安全性

### ✅ 优点
1. **私钥不传输**: 私钥只在 approve-escrow Lambda 内部使用，不通过网络传输
2. **Vault 保护**: 助记词存储在 Vault 中，有访问控制
3. **JWT 认证**: 所有 API 调用都需要有效的 JWT token
4. **一次性 approval**: 使用 MaxUint256，用户只需 approve 一次

### 🔒 安全措施
- approve-escrow Lambda 只能被认证用户调用
- 私钥在内存中使用后立即销毁
- 所有操作都有日志记录
- Vault token 通过环境变量安全传递

## 部署步骤

### 1. 部署 did-login-lambda
```bash
cd did-login-lambda

# 更新 samconfig.toml 添加新参数
# vault_addr = "http://your-vault:8200"
# vault_token = "your-vault-token"
# sepolia_rpc_url = "https://eth-sepolia.g.alchemy.com/v2/YOUR_KEY"

sam build
sam deploy
```

### 2. 部署 xz-wallet-contract lambda
```bash
cd xz-wallet-contract/lambda

# 更新 parameters.json 添加
# {
#   "ParameterKey": "DIDLoginAPIURL",
#   "ParameterValue": "https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod"
# }

sam build
sam deploy
```

## API 接口

### approve-escrow API

**Endpoint**: `POST /api/approve-escrow`

**Headers**:
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Request Body**:
```json
{
  "token_address": "0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8",
  "spender_address": "0x8e98B971884e14C5da6D528932bf96296311B8cb",
  "amount": "optional, defaults to MaxUint256"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "success": true,
    "tx_hash": "0x...",
    "eth_address": "0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA",
    "message": "Escrow contract approved successfully"
  }
}
```

## 使用流程

### 用户视角
1. 用户在 UI 点击"发布任务"
2. 填写任务信息和奖励金额
3. 点击提交
4. 系统自动处理 approval（如果需要）
5. 任务创建成功

### 系统内部流程
1. `create-task` Lambda 收到请求
2. 检查用户的 XZT balance
3. 检查 escrow allowance
4. 如果 allowance 不足:
   - 调用 `approve-escrow` API
   - 等待 approval 交易确认
5. 创建任务到区块链
6. 保存任务到数据库
7. 返回成功响应

## 错误处理

### 常见错误

1. **余额不足**
```json
{
  "success": false,
  "error": "Insufficient XZT balance. Required: 100 XZT, Available: 50 XZT"
}
```

2. **Approval 失败**
```json
{
  "success": false,
  "error": "Failed to approve escrow contract: transaction failed"
}
```

3. **Vault 连接失败**
```json
{
  "success": false,
  "error": "Failed to connect to Vault"
}
```

## 监控和日志

### CloudWatch Logs
- `approve-escrow` Lambda 日志包含:
  - DID 和用户名
  - 派生的以太坊地址
  - 交易哈希
  - 区块确认信息

### 关键指标
- Approval 成功率
- 平均 approval 时间
- Gas 费用统计

## 未来优化

1. **批量 Approval**: 支持一次 approve 多个合约
2. **Gas 优化**: 动态调整 gas price
3. **重试机制**: 交易失败时自动重试
4. **通知系统**: Approval 完成后通知用户

## 相关文件

- `did-login-lambda/go/cmd/approve-escrow/main.go` - Approval Lambda 实现
- `xz-wallet-contract/lambda/cmd/create-task/main.go` - 调用 approval 的逻辑
- `did-login-lambda/template.yaml` - Lambda 配置
- `xz-wallet-contract/lambda/template.yaml` - Lambda 配置
