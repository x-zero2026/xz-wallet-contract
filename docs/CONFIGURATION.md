# XZ Wallet 配置总结

## ✅ 已完成配置

### 1. 智能合约部署 ✓

**网络**: Sepolia Testnet  
**部署时间**: 2026-01-26 10:51:15 UTC

| 合约 | 地址 | 状态 |
|------|------|------|
| XZToken | `0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8` | ✅ 已验证 |
| TaskEscrow | `0x8e98B971884e14C5da6D528932bf96296311B8cb` | ✅ 已验证 |

**Etherscan 链接**:
- [XZToken](https://sepolia.etherscan.io/address/0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8)
- [TaskEscrow](https://sepolia.etherscan.io/address/0x8e98B971884e14C5da6D528932bf96296311B8cb)

---

### 2. 代币分配 ✓

| 账户 | 地址 | XZT 余额 | 用途 |
|------|------|----------|------|
| 系统钱包 | `0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA` | 10,000 XZT | 平台储备金 |
| 部署账户 | `0xd62F159A744df11332F8F1C73C827aed8Ca9378D` | 0 XZT | 已转出 |

**转账交易**: [查看](https://sepolia.etherscan.io/tx/0x2ee03b639dbe0bf132de5d19c38fd6d2dbe10b6f77c71711bad37f593a508199)

---

### 3. 数据库迁移 ✓

**数据库**: Supabase PostgreSQL  
**连接**: `aws-1-ap-south-1.pooler.supabase.com:6543`

**已创建的表**:
- ✅ `users` (扩展: credit_score, xzt_balance, tasks_completed, tasks_cancelled, escrow_approved)
- ✅ `tasks` (任务表)
- ✅ `task_bids` (投标表)
- ✅ `task_submissions` (提交表)
- ✅ `credit_history` (信用分历史)
- ✅ `xzt_transactions` (交易历史)

**已创建的视图**:
- ✅ `v_active_tasks` (活跃任务视图)
- ✅ `v_user_stats` (用户统计视图)

**已创建的触发器**:
- ✅ `trigger_update_tasks_updated_at` (自动更新 tasks.updated_at)
- ✅ `trigger_update_bids_updated_at` (自动更新 task_bids.updated_at)

---

### 4. 用户初始化 ✓

**Admin 用户**:
- Username: `admin`
- XZT Balance: `10,000.00`
- Credit Score: `5,000`
- DID: `0x3070deb1c17432b094d30509ccbfd598fb2793435efdca9273dfbc558bc040ca`
- ETH Address: `0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA`

---

## 🔧 环境变量配置

### 智能合约 (contracts/.env)

```env
# Sepolia Network
SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/GA5ibaTuz122ssPqQhWL7
CHAIN_ID=11155111

# Admin Wallet (Deployer)
ADMIN_WALLET_ADDRESS=0xd62F159A744df11332F8F1C73C827aed8Ca9378D
ADMIN_WALLET_PRIVATE_KEY=<REDACTED>

# Etherscan API Key
ETHERSCAN_API_KEY=VJ63YX9AH7Y9XGE4PSV2D4DSR4B1UKSXYF

# System Wallet
SYSTEM_WALLET_ADDRESS=0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA

# Contract Addresses
XZT_TOKEN_ADDRESS=0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8
TASK_ESCROW_ADDRESS=0x8e98B971884e14C5da6D528932bf96296311B8cb
```

---

### Lambda 函数 (需要配置)

```env
# Database
DATABASE_URL=postgresql://postgres.rbpsksuuvtzmathnmyxn:iPass4xz2026!@aws-1-ap-south-1.pooler.supabase.com:6543/postgres
DB_PASSWORD=iPass4xz2026!

# Blockchain
SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/GA5ibaTuz122ssPqQhWL7
CHAIN_ID=11155111
XZT_TOKEN_ADDRESS=0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8
TASK_ESCROW_ADDRESS=0x8e98B971884e14C5da6D528932bf96296311B8cb

# Admin Wallet (for contract calls)
ADMIN_WALLET_ADDRESS=0xd62F159A744df11332F8F1C73C827aed8Ca9378D
ADMIN_WALLET_PRIVATE_KEY=<REDACTED>

# System Wallet
SYSTEM_WALLET_ADDRESS=0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA

# Vault
VAULT_ADDR=http://your-vault-server:8200
VAULT_TOKEN=your_vault_token_here

# JWT
JWT_SECRET=your_jwt_secret_here
JWT_EXPIRY=168h
```

---

## 📊 系统状态检查

### 检查合约状态

```bash
cd xz-wallet-contract/contracts
npx hardhat console --network sepolia
```

在 console 中执行：

```javascript
// 获取合约实例
const token = await ethers.getContractAt("XZToken", "0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8");
const escrow = await ethers.getContractAt("TaskEscrow", "0x8e98B971884e14C5da6D528932bf96296311B8cb");

// 检查系统钱包余额
const balance = await token.balanceOf("0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA");
console.log("System wallet balance:", ethers.formatEther(balance), "XZT");

// 检查 escrow owner
const owner = await escrow.owner();
console.log("Escrow owner:", owner);

// 检查 escrow token
const tokenAddr = await escrow.token();
console.log("Escrow token:", tokenAddr);
```

---

### 检查数据库状态

```bash
PGPASSWORD='iPass4xz2026!' psql -h aws-1-ap-south-1.pooler.supabase.com -p 6543 -U postgres.rbpsksuuvtzmathnmyxn -d postgres
```

在 psql 中执行：

```sql
-- 检查表
\dt

-- 检查 admin 用户
SELECT username, xzt_balance, credit_score, eth_address 
FROM users 
WHERE username = 'admin';

-- 检查任务表
SELECT COUNT(*) FROM tasks;

-- 检查视图
SELECT * FROM v_user_stats WHERE username = 'admin';
```

---

## 🎯 下一步工作

### 1. 生成 Go 合约绑定 ⏳

```bash
cd xz-wallet-contract/lambda
make generate-bindings
```

这将生成：
- `pkg/blockchain/contracts/xztoken.go`
- `pkg/blockchain/contracts/taskescrow.go`

---

### 2. 实现 Lambda 函数 ⏳

需要实现的函数：

**钱包相关**:
- [ ] `get-balance` - 查询 XZT 余额
- [ ] `transfer-xzt` - 转账 XZT

**任务相关**:
- [ ] `create-task` - 创建任务并锁定 XZT
- [ ] `list-tasks` - 列出任务
- [ ] `get-task` - 获取任务详情
- [ ] `bid-task` - 投标任务
- [ ] `select-bidder` - 选择投标者
- [ ] `submit-work` - 提交工作
- [ ] `approve-work` - 批准工作
- [ ] `reject-work` - 拒绝工作
- [ ] `cancel-task` - 取消任务

**合约交互**:
- [ ] `set-executor` - 设置执行者
- [ ] `pay-milestone` - 支付里程碑

---

### 3. 部署 Lambda 到 AWS ⏳

```bash
cd xz-wallet-contract/lambda
sam build
sam deploy --guided
```

---

### 4. 配置 API Gateway ⏳

- 启用 CORS
- 配置路由
- 设置授权

---

## 📝 快速参考

### 合约地址
```
XZToken:     0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8
TaskEscrow:  0x8e98B971884e14C5da6D528932bf96296311B8cb
```

### 钱包地址
```
Admin DID:      0x3070deb1c17432b094d30509ccbfd598fb2793435efdca9273dfbc558bc040ca
System Wallet:  0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA
Deployer:       0xd62F159A744df11332F8F1C73C827aed8Ca9378D
```

### 数据库
```
Host: aws-1-ap-south-1.pooler.supabase.com
Port: 6543
Database: postgres
User: postgres.rbpsksuuvtzmathnmyxn
```

### Vault
```
Address: http://your-vault-server:8200
Token: your_vault_token_here
```

---

**配置状态**: ✅ 完成  
**准备状态**: ✅ 可以开始后端开发  
**最后更新**: 2026-01-26

