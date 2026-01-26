# Lambda 函数部署指南

## ⚠️ 重要：避免 502 错误

**问题原因：**
使用 `sam deploy` 会更新整个 CloudFormation stack，可能导致某些函数的 bootstrap 文件丢失，造成 502 错误。

**解决方案：**
始终使用以下两步流程部署：

## 正确的部署流程

### 方案 1：更新单个函数（推荐用于小改动）

```bash
# 1. 构建所有函数（确保 bootstrap 文件完整）
sam build

# 2. 只更新修改的函数
cd .aws-sam/build/GetTaskFunction
zip GetTaskFunction.zip bootstrap
cd ../../..
aws lambda update-function-code \
  --function-name xz-wallet-backend-GetTaskFunction-1AGdjuz2dekr \
  --zip-file fileb://.aws-sam/build/GetTaskFunction/GetTaskFunction.zip
```

### 方案 2：批量更新所有函数（推荐用于多个改动）

```bash
# 1. 构建所有函数
sam build

# 2. 运行批量更新脚本
./update-all-functions.sh
```

### 方案 3：新增函数时使用 sam deploy

**只在以下情况使用 `sam deploy`：**
- 新增 Lambda 函数
- 修改 API Gateway 路由
- 修改 IAM 权限
- 修改环境变量

**步骤：**
```bash
# 1. 构建所有函数
sam build

# 2. 部署（会更新整个 stack）
sam deploy

# 3. 如果出现 502，立即运行修复脚本
./update-all-functions.sh
```

## 当前函数列表

| 函数名 | AWS 函数名 | API 路径 |
|--------|-----------|----------|
| GetBalanceFunction | xz-wallet-backend-GetBalanceFunction-lmyirQsvhGBD | GET /wallet/balance |
| CreateTaskFunction | xz-wallet-backend-CreateTaskFunction-iJeNYr7lqfWS | POST /tasks |
| ListTasksFunction | xz-wallet-backend-ListTasksFunction-UsWjqkcZmMZC | GET /tasks |
| GetTaskFunction | xz-wallet-backend-GetTaskFunction-1AGdjuz2dekr | GET /tasks/{id} |
| BidTaskFunction | xz-wallet-backend-BidTaskFunction-kZY7RrJ8UT3k | POST /tasks/{id}/bid |
| SelectBidderFunction | xz-wallet-backend-SelectBidderFunction-p8xCW1224ODx | POST /tasks/{id}/select-bidder |
| ApproveWorkFunction | xz-wallet-backend-ApproveWorkFunction-mAoNJHBZeELU | POST /tasks/{id}/approve |
| CancelTaskFunction | xz-wallet-backend-CancelTaskFunction-OPpwD7mkgToS | POST /tasks/{id}/cancel |
| SubmitWorkFunction | xz-wallet-backend-SubmitWorkFunction-iJFNYYaUx7Fj | POST /tasks/{id}/submit |

## 快速命令参考

### 检查所有函数状态
```bash
aws lambda list-functions \
  --query "Functions[?starts_with(FunctionName, 'xz-wallet-backend-')].[FunctionName,LastUpdateStatus]" \
  --output table
```

### 查看单个函数状态
```bash
aws lambda get-function \
  --function-name xz-wallet-backend-GetTaskFunction-1AGdjuz2dekr \
  --query 'Configuration.LastUpdateStatus'
```

### 查看函数日志
```bash
aws logs tail /aws/lambda/xz-wallet-backend-GetTaskFunction-1AGdjuz2dekr --follow
```

## 故障排除

### 如果遇到 502 错误

1. **立即运行修复脚本：**
   ```bash
   sam build
   ./update-all-functions.sh
   ```

2. **检查函数状态：**
   ```bash
   aws lambda list-functions \
     --query "Functions[?starts_with(FunctionName, 'xz-wallet-backend-')].[FunctionName,State,LastUpdateStatus]" \
     --output table
   ```

3. **查看错误日志：**
   ```bash
   aws logs tail /aws/lambda/xz-wallet-backend-FUNCTION-NAME --follow
   ```

### 常见错误

| 错误 | 原因 | 解决方案 |
|------|------|----------|
| 502 Bad Gateway | bootstrap 文件丢失 | 运行 `update-all-functions.sh` |
| 403 Forbidden | CORS 配置问题 | 检查 template.yaml 中的 CORS 设置 |
| 401 Unauthorized | JWT token 问题 | 检查前端 token 是否正确传递 |
| 500 Internal Server Error | 代码逻辑错误 | 查看 CloudWatch 日志 |

## 最佳实践

1. ✅ **每次修改代码后先 `sam build`**
2. ✅ **优先使用 `update-all-functions.sh` 而不是 `sam deploy`**
3. ✅ **新增函数后立即运行 `update-all-functions.sh` 修复**
4. ✅ **部署前备份当前工作的函数代码**
5. ✅ **在本地测试后再部署到 AWS**

## 开发工作流

```bash
# 1. 修改代码
vim cmd/get-task/main.go

# 2. 本地测试（可选）
sam build
sam local invoke GetTaskFunction -e events/get-task.json

# 3. 部署到 AWS
sam build
./update-all-functions.sh

# 4. 验证
curl https://yms07x0sn0.execute-api.us-east-1.amazonaws.com/prod/tasks/xxx
```

## 紧急回滚

如果部署出现严重问题：

```bash
# 1. 回滚到上一个 CloudFormation 版本
aws cloudformation cancel-update-stack --stack-name xz-wallet-backend

# 2. 或者重新部署上一个工作版本
git checkout <last-working-commit>
sam build
./update-all-functions.sh
```
