# XZ Wallet & Task Escrow System - MVP Requirements

## 📋 MVP Overview

**Version**: MVP 1.0 (Simplified)  
**Goal**: Validate core functionality with minimal complexity  
**Timeline**: 1-2 weeks  

### MVP Scope

**✅ Included**:
- Wallet management (Ethereum-compatible DID)
- XZT token on Sepolia (real blockchain)
- Task creation with XZT escrow
- Bidding mechanism with credit score
- Milestone-based payments (30%, 80%, 100%)
- Task cancellation with refunds
- Credit score system

**❌ Excluded (Future)**:
- Signature verification (simplified to admin-only)
- Complex state machine in smart contract
- Event listening and sync
- Faucet system
- Transfer limits enforcement in contract

### Key Simplifications

| Feature | Full Version | MVP Version |
|---------|-------------|-------------|
| State Machine | On-chain in contract | In database only |
| Payment Trigger | Contract auto-payment | Backend calls contract |
| User Gas | User pays (needs ETH) | Backend pays (no ETH needed) |
| Signature | User signs each action | Admin-only contract calls |
| Approval | Per-action approval | One-time unlimited approval |

---

## 🎯 Core Features (MVP)

### 1. Wallet Management

**User Wallet**:
- DID = Ethereum address (secp256k1)
- Mnemonic stored in Vault
- Shown once during registration
- No export/import
- **No Sepolia ETH required** (backend pays gas)

**System Wallet**:
- Holds initial 10,000 XZT
- Used for manual distribution

**Admin Wallet**:
- Your MetaMask wallet
- Pays all gas fees
- Calls contracts on behalf of users

### 2. XZT Token (Real Blockchain)

- Standard ERC20 on Sepolia
- Initial mint: 10,000 XZT to system wallet
- Users can receive and hold XZT
- Users can transfer XZT (direct ERC20 transfer)

### 3. Task Escrow (Hybrid: Database + Blockchain)

**Database Stores**:
- Task metadata (name, description, criteria)
- Task status (state machine)
- Bid applications
- Submissions
- Credit scores

**Blockchain Stores**:
- Locked XZT amounts
- Payment records
- Task-executor mapping

**Payment Flow**:
```
Creator creates task 
→ Backend locks XZT in contract (transferFrom)
→ Database records task as PENDING

Executor completes milestone
→ Creator approves in UI
→ Backend updates database status
→ Backend calls contract to pay milestone
→ XZT transferred to executor
```

---

## 🔄 Task State Machine (Database Only)

### States
```
PENDING              → Task published, waiting for bids
BIDDING              → Users applied, creator selecting
ACCEPTED             → Executor selected
DESIGN_SUBMITTED     → Design submitted
DESIGN_APPROVED      → Design approved, 30% paid
IMPLEMENTATION_SUBMITTED → Implementation submitted
IMPLEMENTATION_APPROVED  → Implementation approved, 80% paid
FINAL_SUBMITTED      → Final work submitted
COMPLETED            → Final approved, 100% paid
CANCELLED            → Task cancelled
```

### State Transitions (Backend Controlled)

All state transitions happen in database via backend APIs. No on-chain state machine.

---

## 💰 Payment Rules (Same as Full Version)

### Milestone Payments
| Milestone | Payment | Cumulative | Trigger |
|-----------|---------|------------|---------|
| Design Approved | 30% | 30% | Backend calls `payMilestone(taskId, 30%)` |
| Implementation Approved | 50% | 80% | Backend calls `payMilestone(taskId, 50%)` |
| Final Approved | 20% | 100% | Backend calls `payMilestone(taskId, 20%)` |

### Cancellation Payments (Same as Full Version)

**Creator Cancels**:
| Current State | Executor Gets | Creator Gets Back |
|---------------|---------------|-------------------|
| Before DESIGN_APPROVED | 0% | 100% |
| After DESIGN_APPROVED | 30% | 70% |
| After IMPLEMENTATION_APPROVED | 80% | 20% |

**Executor Cancels**:
| Current State | Executor Keeps | Creator Gets Back | Credit Penalty |
|---------------|----------------|-------------------|----------------|
| Before DESIGN_APPROVED | 0% | 100% | 0 |
| After DESIGN_APPROVED | 30% | 70% | -3,000 |
| After IMPLEMENTATION_APPROVED | 80% | 20% | -8,000 |

---

## 🏆 Credit Score System (Same as Full Version)

### Initial Score
- Every user: **5,000 points**

### Score Changes
- Cancel after design approved: **-3,000**
- Cancel after implementation approved: **-8,000**
- Task completed: **+100**

### Rules
- Credit < 0: Cannot bid
- Credit < 1,000: Warning badge
- No recovery except completing tasks

---

## 🔐 Security (Simplified for MVP)

### Key Management
- User mnemonics: Vault
- Admin private key: Environment variable (your MetaMask key)
- System wallet mnemonic: Vault

### Contract Access Control
- Only admin wallet can call contract functions
- No signature verification (simplified)
- Backend validates user permissions before calling contract

### User Authorization Flow
```
User clicks "Approve Design"
→ Frontend calls backend API with JWT
→ Backend validates:
  - User is task creator
  - Task status is DESIGN_SUBMITTED
→ Backend updates database
→ Backend calls contract: payMilestone(taskId, 30%)
→ Contract transfers XZT to executor
```

---

## 🏗️ Smart Contracts (Simplified)

### XZToken.sol (Standard ERC20)
```solidity
contract XZToken is ERC20 {
    constructor() ERC20("XZ Token", "XZT") {
        _mint(msg.sender, 10000 * 10**18);
    }
}
```

### TaskEscrow.sol (Simplified)
```solidity
contract TaskEscrow {
    IERC20 public token;
    address public admin;
    
    struct Task {
        address creator;
        address executor;
        uint256 totalAmount;
        uint256 paidAmount;
    }
    
    mapping(uint256 => Task) public tasks;
    uint256 public nextTaskId;
    
    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin");
        _;
    }
    
    // Create task and lock XZT
    function createTask(
        address creator,
        address executor,
        uint256 amount
    ) external onlyAdmin returns (uint256) {
        token.transferFrom(creator, address(this), amount);
        
        uint256 taskId = nextTaskId++;
        tasks[taskId] = Task(creator, executor, amount, 0);
        
        emit TaskCreated(taskId, creator, executor, amount);
        return taskId;
    }
    
    // Pay milestone to executor
    function payMilestone(
        uint256 taskId,
        uint256 amount
    ) external onlyAdmin {
        Task storage task = tasks[taskId];
        require(task.paidAmount + amount <= task.totalAmount, "Exceeds total");
        
        token.transfer(task.executor, amount);
        task.paidAmount += amount;
        
        emit MilestonePaid(taskId, task.executor, amount);
    }
    
    // Cancel task with refund
    function cancelTask(
        uint256 taskId,
        uint256 executorAmount
    ) external onlyAdmin {
        Task storage task = tasks[taskId];
        uint256 remaining = task.totalAmount - task.paidAmount;
        
        if (executorAmount > 0) {
            token.transfer(task.executor, executorAmount);
        }
        
        uint256 refund = remaining - executorAmount;
        if (refund > 0) {
            token.transfer(task.creator, refund);
        }
        
        task.paidAmount = task.totalAmount;
        
        emit TaskCancelled(taskId, executorAmount, refund);
    }
    
    // Update executor (for bidding)
    function setExecutor(
        uint256 taskId,
        address newExecutor
    ) external onlyAdmin {
        tasks[taskId].executor = newExecutor;
        emit ExecutorSet(taskId, newExecutor);
    }
    
    event TaskCreated(uint256 indexed taskId, address creator, address executor, uint256 amount);
    event MilestonePaid(uint256 indexed taskId, address executor, uint256 amount);
    event TaskCancelled(uint256 indexed taskId, uint256 executorAmount, uint256 refund);
    event ExecutorSet(uint256 indexed taskId, address executor);
}
```

---

## 📊 Database Schema (MVP)

### Extended Users Table
```sql
ALTER TABLE users 
    ADD COLUMN credit_score INT DEFAULT 5000,
    ADD COLUMN tasks_completed INT DEFAULT 0,
    ADD COLUMN tasks_cancelled INT DEFAULT 0,
    ADD COLUMN xzt_balance DECIMAL(20, 8) DEFAULT 0;
```

### Tasks Table
```sql
CREATE TABLE tasks (
    task_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_task_id BIGINT UNIQUE,
    project_id UUID NOT NULL REFERENCES projects(project_id),
    creator_did VARCHAR(42) NOT NULL REFERENCES users(did),
    executor_did VARCHAR(42) REFERENCES users(did),
    
    task_name VARCHAR(255) NOT NULL,
    task_description TEXT NOT NULL,
    acceptance_criteria TEXT NOT NULL,
    reward_amount DECIMAL(20, 8) NOT NULL,
    paid_amount DECIMAL(20, 8) DEFAULT 0,
    
    visibility VARCHAR(20) NOT NULL CHECK (visibility IN ('project', 'global')),
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_creator ON tasks(creator_did);
CREATE INDEX idx_tasks_executor ON tasks(executor_did);
CREATE INDEX idx_tasks_project ON tasks(project_id);
```

### Task Bids Table
```sql
CREATE TABLE task_bids (
    bid_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(task_id) ON DELETE CASCADE,
    bidder_did VARCHAR(42) NOT NULL REFERENCES users(did),
    bid_message TEXT,
    credit_score_snapshot INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(task_id, bidder_did)
);

CREATE INDEX idx_bids_task ON task_bids(task_id);
CREATE INDEX idx_bids_bidder ON task_bids(bidder_did);
```

### Task Submissions Table
```sql
CREATE TABLE task_submissions (
    submission_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(task_id) ON DELETE CASCADE,
    submission_type VARCHAR(30) NOT NULL,
    content TEXT NOT NULL,
    file_urls TEXT[],
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_submissions_task ON task_submissions(task_id);
```

### Credit History Table
```sql
CREATE TABLE credit_history (
    history_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_did VARCHAR(42) NOT NULL REFERENCES users(did),
    task_id UUID REFERENCES tasks(task_id),
    change_amount INT NOT NULL,
    reason VARCHAR(50) NOT NULL,
    before_score INT NOT NULL,
    after_score INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_credit_user ON credit_history(user_did, created_at DESC);
```

---

## 🔌 API Endpoints (MVP)

### Wallet APIs

**GET /wallet/balance**
```json
Response:
{
  "xzt_balance": "1000.00",
  "did": "0x..."
}
```

**POST /wallet/transfer**
```json
Request:
{
  "to_address": "0x...",
  "amount": "100.00"
}

Response:
{
  "tx_hash": "0x...",
  "status": "confirmed"
}
```

Note: Direct ERC20 transfer, backend calls `token.transfer()` on behalf of user.

---

### Task APIs

**POST /tasks**
```json
Request:
{
  "project_id": "uuid",
  "task_name": "Build login page",
  "task_description": "...",
  "acceptance_criteria": "...",
  "reward_amount": "5000.00",
  "visibility": "project"
}

Process:
1. Validate user balance >= reward_amount
2. Call contract: createTask(creator, 0x0, amount)
3. Save to database with contract_task_id
4. Return task_id

Response:
{
  "task_id": "uuid",
  "contract_task_id": 123,
  "status": "pending"
}
```

**GET /tasks**
```json
Query: ?visibility=project&project_id=uuid&status=pending

Response:
{
  "tasks": [
    {
      "task_id": "uuid",
      "task_name": "...",
      "reward_amount": "5000.00",
      "paid_amount": "0.00",
      "status": "pending",
      "creator": {
        "did": "0x...",
        "username": "alice"
      },
      "bid_count": 5
    }
  ]
}
```

**GET /tasks/:id**
```json
Response:
{
  "task_id": "uuid",
  "contract_task_id": 123,
  "task_name": "...",
  "reward_amount": "5000.00",
  "paid_amount": "1500.00",
  "status": "design_approved",
  "creator": {...},
  "executor": {...},
  "submissions": [...]
}
```

**POST /tasks/:id/bid**
```json
Request:
{
  "message": "I have experience..."
}

Process:
1. Validate credit_score >= 0
2. Save bid to database
3. Update task status to BIDDING if first bid

Response:
{
  "bid_id": "uuid"
}
```

**GET /tasks/:id/bids** (Creator only)
```json
Response:
{
  "bids": [
    {
      "bidder": {
        "did": "0x...",
        "username": "bob",
        "credit_score": 4500,
        "tasks_completed": 10,
        "tasks_cancelled": 1
      },
      "message": "..."
    }
  ]
}
```

**POST /tasks/:id/select-bidder** (Creator only)
```json
Request:
{
  "bidder_did": "0x..."
}

Process:
1. Validate user is creator
2. Call contract: setExecutor(taskId, bidderDid)
3. Update database: executor_did, status=ACCEPTED
4. Reject other bids

Response:
{
  "status": "accepted"
}
```

**POST /tasks/:id/submit** (Executor only)
```json
Request:
{
  "submission_type": "design",
  "content": "...",
  "file_urls": ["https://..."]
}

Process:
1. Validate user is executor
2. Save submission to database
3. Update task status (ACCEPTED→DESIGN_SUBMITTED, etc.)

Response:
{
  "submission_id": "uuid",
  "status": "design_submitted"
}
```

**POST /tasks/:id/approve** (Creator only)
```json
Request:
{
  "milestone": "design"
}

Process:
1. Validate user is creator
2. Calculate payment amount (30%, 50%, or 20%)
3. Call contract: payMilestone(taskId, amount)
4. Update database: status, paid_amount
5. If final milestone, update credit_score (+100)

Response:
{
  "status": "design_approved",
  "payment": {
    "amount": "1500.00",
    "tx_hash": "0x..."
  }
}
```

**POST /tasks/:id/reject** (Creator only)
```json
Request:
{
  "milestone": "design",
  "reason": "..."
}

Process:
1. Validate user is creator
2. Revert status (DESIGN_SUBMITTED→ACCEPTED, etc.)

Response:
{
  "status": "accepted"
}
```

**POST /tasks/:id/cancel**
```json
Request:
{
  "reason": "..."
}

Process:
1. Validate user is creator or executor
2. Calculate refund based on current status
3. Call contract: cancelTask(taskId, executorAmount)
4. Update database: status=CANCELLED
5. Update credit_score if executor cancels

Response:
{
  "status": "cancelled",
  "refund": {
    "creator_refund": "3500.00",
    "executor_payment": "1500.00"
  },
  "credit_penalty": -3000
}
```

---

## 🚀 Implementation Plan (MVP)

### Week 1: Foundation

**Day 1-2: Smart Contracts**
- [ ] Write XZToken.sol
- [ ] Write TaskEscrow.sol (simplified)
- [ ] Write deployment script
- [ ] Deploy to Sepolia
- [ ] Mint 10,000 XZT to system wallet
- [ ] Verify contracts on Etherscan

**Day 3-4: DID & Database**
- [ ] Modify DID generation (secp256k1)
- [ ] Create database migration script
- [ ] Run migrations on Supabase
- [ ] Test DID generation

**Day 5: Contract Integration**
- [ ] Create Go package for contract interaction
- [ ] Implement contract call functions
- [ ] Test contract calls on Sepolia

---

### Week 2: Backend & Frontend

**Day 1-2: Backend APIs**
- [ ] Wallet APIs (balance, transfer)
- [ ] Task CRUD APIs
- [ ] Bidding APIs
- [ ] Submission APIs
- [ ] Approval/Rejection APIs
- [ ] Cancellation APIs

**Day 3-4: Frontend UI**
- [ ] Wallet page (balance, transfer)
- [ ] Task list page
- [ ] Task detail page
- [ ] Bid interface
- [ ] Submission forms
- [ ] Approval buttons

**Day 5: Testing & Deployment**
- [ ] Integration testing
- [ ] Deploy backend to Lambda
- [ ] Deploy frontend to Amplify
- [ ] End-to-end testing

---

## 🔧 Technical Implementation Details

### DID Generation (Updated)

```go
package auth

import (
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/tyler-smith/go-bip39"
    "github.com/tyler-smith/go-bip32"
)

// GenerateDIDAndMnemonic generates Ethereum-compatible DID
func GenerateDIDAndMnemonic() (did string, mnemonic string, err error) {
    // Generate mnemonic
    entropy, _ := bip39.NewEntropy(128)
    mnemonic, _ = bip39.NewMnemonic(entropy)
    
    // Generate seed
    seed := bip39.NewSeed(mnemonic, "")
    
    // Derive key using BIP44 path: m/44'/60'/0'/0/0
    masterKey, _ := bip32.NewMasterKey(seed)
    purposeKey, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
    coinKey, _ := purposeKey.NewChildKey(bip32.FirstHardenedChild + 60)
    accountKey, _ := coinKey.NewChildKey(bip32.FirstHardenedChild + 0)
    changeKey, _ := accountKey.NewChildKey(0)
    addressKey, _ := changeKey.NewChildKey(0)
    
    // Get private key
    privateKey, _ := crypto.ToECDSA(addressKey.Key)
    
    // Get Ethereum address
    address := crypto.PubkeyToAddress(privateKey.PublicKey)
    did = address.Hex() // 0x...
    
    return did, mnemonic, nil
}

// MnemonicToDID converts mnemonic to DID
func MnemonicToDID(mnemonic string) (string, error) {
    seed := bip39.NewSeed(mnemonic, "")
    // ... same derivation as above
    return did, nil
}

// MnemonicToPrivateKey for signing
func MnemonicToPrivateKey(mnemonic string) (*ecdsa.PrivateKey, error) {
    seed := bip39.NewSeed(mnemonic, "")
    // ... same derivation
    return privateKey, nil
}
```

### Contract Interaction

```go
package blockchain

import (
    "context"
    "math/big"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/ethclient"
)

type ContractClient struct {
    client       *ethclient.Client
    token        *XZToken
    escrow       *TaskEscrow
    adminAuth    *bind.TransactOpts
}

// CreateTask locks XZT in escrow
func (c *ContractClient) CreateTask(
    creatorAddress string,
    amount *big.Int,
) (uint64, error) {
    // First, approve escrow to spend creator's XZT
    // (This should be done once per user, not per task)
    
    // Create task
    tx, err := c.escrow.CreateTask(
        c.adminAuth,
        common.HexToAddress(creatorAddress),
        common.HexToAddress("0x0"), // executor TBD
        amount,
    )
    if err != nil {
        return 0, err
    }
    
    // Wait for confirmation
    receipt, err := bind.WaitMined(context.Background(), c.client, tx)
    if err != nil {
        return 0, err
    }
    
    // Parse TaskCreated event to get taskId
    taskId := parseTaskIdFromReceipt(receipt)
    
    return taskId, nil
}

// PayMilestone pays executor
func (c *ContractClient) PayMilestone(
    taskId uint64,
    amount *big.Int,
) (string, error) {
    tx, err := c.escrow.PayMilestone(
        c.adminAuth,
        big.NewInt(int64(taskId)),
        amount,
    )
    if err != nil {
        return "", err
    }
    
    receipt, err := bind.WaitMined(context.Background(), c.client, tx)
    if err != nil {
        return "", err
    }
    
    return tx.Hash().Hex(), nil
}

// CancelTask with refund
func (c *ContractClient) CancelTask(
    taskId uint64,
    executorAmount *big.Int,
) (string, error) {
    tx, err := c.escrow.CancelTask(
        c.adminAuth,
        big.NewInt(int64(taskId)),
        executorAmount,
    )
    if err != nil {
        return "", err
    }
    
    receipt, err := bind.WaitMined(context.Background(), c.client, tx)
    if err != nil {
        return "", err
    }
    
    return tx.Hash().Hex(), nil
}

// SetExecutor updates executor address
func (c *ContractClient) SetExecutor(
    taskId uint64,
    executorAddress string,
) error {
    tx, err := c.escrow.SetExecutor(
        c.adminAuth,
        big.NewInt(int64(taskId)),
        common.HexToAddress(executorAddress),
    )
    if err != nil {
        return err
    }
    
    _, err = bind.WaitMined(context.Background(), c.client, tx)
    return err
}
```

### One-Time Approval Setup

```go
// ApproveEscrowForUser - called once per user
func (c *ContractClient) ApproveEscrowForUser(
    userAddress string,
    userMnemonic string,
) error {
    // Get user's private key
    userPrivateKey, _ := auth.MnemonicToPrivateKey(userMnemonic)
    
    // Create user's auth
    userAuth, _ := bind.NewKeyedTransactorWithChainID(
        userPrivateKey,
        big.NewInt(11155111), // Sepolia chain ID
    )
    
    // But use admin's gas settings
    userAuth.GasLimit = 100000
    userAuth.GasPrice = c.adminAuth.GasPrice
    
    // Approve unlimited amount
    maxUint256 := new(big.Int)
    maxUint256.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
    
    tx, err := c.token.Approve(
        userAuth,
        c.escrow.Address,
        maxUint256,
    )
    if err != nil {
        return err
    }
    
    _, err = bind.WaitMined(context.Background(), c.client, tx)
    return err
}
```

Note: The above approach requires user to have ETH for gas. For MVP, we can:
1. Give each new user 0.01 Sepolia ETH from admin wallet
2. Or implement a meta-transaction pattern (more complex)

**Recommended for MVP**: Give 0.01 ETH to new users once.

---

## 📝 Environment Variables

```bash
# Database
DATABASE_URL=postgresql://postgres:iPass4xz2026!@aws-1-ap-south-1.pooler.supabase.com:6543/postgres

# Blockchain
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
CHAIN_ID=11155111
XZT_TOKEN_ADDRESS=0x...  # After deployment
TASK_ESCROW_ADDRESS=0x...  # After deployment

# Vault
VAULT_ADDR=https://vault.example.com
VAULT_TOKEN=...
VAULT_MNEMONIC_PATH=secret/data/xz-wallet/mnemonics

# Admin Wallet (Your MetaMask)
ADMIN_WALLET_ADDRESS=0x...
ADMIN_WALLET_PRIVATE_KEY=0x...

# System Wallet (Generated)
SYSTEM_WALLET_ADDRESS=0x...

# JWT
JWT_SECRET=your_jwt_secret_here

# API
API_BASE_URL=https://ynnid7kam5.execute-api.us-east-1.amazonaws.com/prod
```

---

## ⚠️ MVP Limitations & Future Improvements

### Current Limitations

1. **No Signature Verification**
   - Current: Admin-only contract calls
   - Future: User signature verification

2. **State Machine in Database**
   - Current: All states in database
   - Future: Critical states on-chain

3. **Centralized Payment Trigger**
   - Current: Backend calls contract to pay
   - Future: Contract auto-pays on approval

4. **Gas Management**
   - Current: Admin pays all gas
   - Future: Users pay their own gas or meta-transactions

5. **No Event Listening**
   - Current: No sync between chain and database
   - Future: Listen to contract events and sync

### Migration Path to Full Version

**Phase 1 (MVP)**: Database state + Contract payments ✅  
**Phase 2**: Add signature verification  
**Phase 3**: Move state machine to contract  
**Phase 4**: Implement event listening  
**Phase 5**: Add meta-transactions for gas-less UX  

---

## ✅ MVP Acceptance Criteria

### Must Have
- ✅ Users can register with Ethereum DID
- ✅ Users can view XZT balance
- ✅ Users can transfer XZT
- ✅ Users can create tasks (XZT locked on-chain)
- ✅ Users can bid on tasks
- ✅ Creators can select bidders
- ✅ Executors can submit work
- ✅ Creators can approve/reject
- ✅ Automatic payment on approval (30%, 80%, 100%)
- ✅ Task cancellation with refunds
- ✅ Credit score tracking

### Nice to Have (Can Skip for MVP)
- ⏭️ Transfer limits enforcement
- ⏭️ Faucet system
- ⏭️ Transaction history page
- ⏭️ Credit score history page
- ⏭️ Task search and filters
- ⏭️ Email notifications

### Performance
- API response < 2s (excluding blockchain confirmation)
- Blockchain confirmation: 15-30s (acceptable for MVP)

---

## 🎯 Success Metrics

### Technical
- [ ] All contracts deployed and verified on Sepolia
- [ ] 10,000 XZT minted to system wallet
- [ ] All APIs functional and tested
- [ ] Frontend deployed and accessible

### Functional
- [ ] Complete one full task workflow end-to-end
- [ ] Test all payment scenarios (30%, 80%, 100%)
- [ ] Test cancellation scenarios
- [ ] Verify credit score updates correctly

### User Experience
- [ ] User can complete task workflow without errors
- [ ] Payment reflects in wallet within 30 seconds
- [ ] UI is responsive and intuitive

---

## 📚 Next Steps After MVP

1. **Gather Feedback**: Test with real users
2. **Add Signature Verification**: Improve security
3. **Optimize Gas**: Reduce transaction costs
4. **Add Features**: Transfer limits, faucet, notifications
5. **Scale**: Handle more concurrent users
6. **Migrate to Mainnet**: Deploy to Ethereum mainnet

---

**Document Version**: MVP 1.0  
**Last Updated**: 2026-01-26  
**Status**: Ready for Implementation  
**Estimated Timeline**: 1-2 weeks
