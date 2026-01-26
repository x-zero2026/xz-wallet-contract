# XZ Wallet & Task Escrow System - Requirements Document

## 📋 Project Overview

**Project Name**: XZ Wallet & Task Escrow System  
**Blockchain**: Sepolia Testnet  
**Token**: XZT (ERC20) - 1 XZT = 1 e-CNY  
**Backend**: Go + AWS Lambda  
**Database**: PostgreSQL (Supabase)  
**Key Management**: HashiCorp Vault  

---

## 🎯 Core Features

### 1. Wallet Management
- User wallet = DID (Ethereum address derived from mnemonic)
- Each user has independent mnemonic stored in Vault
- Mnemonic shown once during registration (for backup)
- No export/import functionality
- Backend manages private keys for signing

### 2. Token System
- **XZT Token** (ERC20)
- Initial mint: 10,000 XZT to system wallet
- 1 XZT = 1 e-CNY (pegged to RMB)

### 3. Transfer Functionality
- Query wallet balance (XZT + Sepolia ETH)
- Transfer XZT to other users
- Transfer limits (backend enforced):
  - Individual: max 1,000 XZT per transaction
  - System wallet: max 30% of total balance per transaction
  - Limits only apply to regular transfers, not task payments

### 4. Task Escrow System
- Milestone-based payment (30%, 80%, 100%)
- Bidding mechanism with credit score
- Decentralized contract execution
- Automatic payment on milestone approval

---

## 🔄 Task State Machine

### States
```
PENDING              → Task published, waiting for bids
BIDDING              → Users applied, creator selecting
ACCEPTED             → Executor selected, starting design
DESIGN_SUBMITTED     → Design submitted, waiting for approval
DESIGN_APPROVED      → Design approved, 30% paid
IMPLEMENTATION_SUBMITTED → Implementation submitted, waiting for approval
IMPLEMENTATION_APPROVED  → Implementation approved, 80% paid (cumulative)
FINAL_SUBMITTED      → Final work submitted, waiting for approval
COMPLETED            → Final approved, 100% paid (cumulative)
CANCELLED            → Task cancelled
```

### State Transitions

**1. Publishing → Bidding**
```
PENDING → (users apply) → BIDDING
```

**2. Bidding → Accepted**
```
BIDDING → (creator selects executor) → ACCEPTED
- Other bids automatically rejected
- No more bids allowed
```

**3. Design Phase**
```
ACCEPTED → (executor submits design) → DESIGN_SUBMITTED
DESIGN_SUBMITTED → (creator approves) → DESIGN_APPROVED [Pay 30%]
DESIGN_SUBMITTED → (creator rejects) → ACCEPTED [Resubmit, unlimited times]
```

**4. Implementation Phase**
```
DESIGN_APPROVED → (executor submits implementation) → IMPLEMENTATION_SUBMITTED
IMPLEMENTATION_SUBMITTED → (creator approves) → IMPLEMENTATION_APPROVED [Pay 50%, total 80%]
IMPLEMENTATION_SUBMITTED → (creator rejects) → DESIGN_APPROVED [Resubmit, unlimited times]
```

**5. Final Phase**
```
IMPLEMENTATION_APPROVED → (executor submits final) → FINAL_SUBMITTED
FINAL_SUBMITTED → (creator approves) → COMPLETED [Pay 20%, total 100%]
FINAL_SUBMITTED → (creator rejects) → IMPLEMENTATION_APPROVED [Resubmit, unlimited times]
```

**6. Cancellation**
```
Any state → (creator/executor cancels) → CANCELLED
```

---

## 💰 Payment Rules

### Milestone Payments (Automatic by Contract)
| Milestone | Payment | Cumulative | Trigger |
|-----------|---------|------------|---------|
| Design Approved | 30% | 30% | Creator approves design |
| Implementation Approved | 50% | 80% | Creator approves implementation |
| Final Approved | 20% | 100% | Creator approves final work |

### Cancellation Payments

**Creator Cancels:**
| Current State | Executor Gets | Creator Gets Back |
|---------------|---------------|-------------------|
| PENDING / BIDDING / ACCEPTED / DESIGN_SUBMITTED | 0% | 100% |
| DESIGN_APPROVED / IMPLEMENTATION_SUBMITTED | 30% (already paid) | 70% |
| IMPLEMENTATION_APPROVED / FINAL_SUBMITTED | 80% (already paid) | 20% |

**Executor Cancels:**
| Current State | Executor Keeps | Creator Gets Back | Credit Score Penalty |
|---------------|----------------|-------------------|---------------------|
| ACCEPTED / DESIGN_SUBMITTED | 0% | 100% | 0 |
| DESIGN_APPROVED / IMPLEMENTATION_SUBMITTED | 30% (already paid) | 70% | -3,000 |
| IMPLEMENTATION_APPROVED / FINAL_SUBMITTED | 80% (already paid) | 20% | -8,000 |

---

## 🏆 Credit Score System

### Initial Score
- Every user starts with: **5,000 credit score**

### Score Changes

**Penalties (Executor Cancels):**
- Cancel after design approved: **-3,000 points**
- Cancel after implementation approved: **-8,000 points**
- Cancel before design approved: **0 points**

**Rewards (Task Completed):**
- Task successfully completed: **+100 points**

### Credit Score Rules
- Credit score < 0: **Cannot bid on tasks**
- Credit score < 1,000: **Warning badge shown to creators**
- No recovery mechanism (cannot buy back credit score)
- Only way to recover: complete tasks successfully

### Bidding Display
Creators can see for each bidder:
- ✅ Credit score
- ✅ Total tasks completed
- ✅ Total tasks cancelled
- ❌ No rating system (not implemented)

---

## 🔐 Security & Signing

### Signature Mechanism

**Purpose**: Verify user authorization for contract operations

**Signing Process**:
1. User clicks action in frontend (e.g., "Approve Design")
2. Frontend calls backend API
3. Backend retrieves user's mnemonic from Vault
4. Backend derives private key and generates signature
5. Backend calls smart contract with signature
6. Contract verifies signature matches expected user address

**Signature Format**:
```solidity
message = keccak256(abi.encodePacked(
    taskId,
    action,      // "approveDesign", "approveImplementation", etc.
    deadline     // Unix timestamp (current time + 5 minutes)
))

signature = sign(message, userPrivateKey)
```

**Contract Verification**:
```solidity
function verifySignature(
    uint256 taskId,
    string memory action,
    uint256 deadline,
    bytes memory signature
) internal view returns (address) {
    require(block.timestamp <= deadline, "Signature expired");
    
    bytes32 message = keccak256(abi.encodePacked(taskId, action, deadline));
    bytes32 ethSignedMessage = keccak256(abi.encodePacked(
        "\x19Ethereum Signed Message:\n32",
        message
    ));
    
    return ecrecover(ethSignedMessage, v, r, s);
}
```

### Key Management

**User Mnemonics**:
- Storage: HashiCorp Vault
- Access: Backend only
- Shown to user: Once during registration (for backup)
- Export: Not allowed
- Import: Not allowed

**Admin Wallet**:
- Purpose: Pay gas fees for contract calls
- Storage: User's MetaMask (not generated by system)
- Private key: Not stored in system

**System Wallet**:
- Purpose: Hold initial 10,000 XZT for distribution
- Storage: Mnemonic in Vault
- Will be generated during deployment

---

## 🏗️ Technical Architecture

### Smart Contracts

**1. XZToken.sol (ERC20)**
```solidity
contract XZToken is ERC20 {
    - Standard ERC20 functionality
    - Initial mint: 10,000 XZT to system wallet
    - No transfer restrictions in contract
}
```

**2. TaskEscrow.sol**
```solidity
contract TaskEscrow {
    - Task creation (lock XZT in contract)
    - Bidding mechanism (off-chain, only selected executor recorded)
    - Milestone approvals with signature verification
    - Automatic payments on approval
    - Cancellation with refund logic
    - All state transitions on-chain
}
```

### Backend Architecture

**Technology Stack**:
- Language: Go
- Deployment: AWS Lambda
- API Gateway: AWS API Gateway
- Database: PostgreSQL (Supabase)
- Key Storage: HashiCorp Vault
- Blockchain: Sepolia Testnet

**API Flow**:
```
User Action (Frontend)
  ↓
Backend API (Lambda)
  ↓
├─ Read mnemonic from Vault
├─ Derive private key
├─ Generate signature
├─ Call smart contract (wait for confirmation)
├─ Update database
  ↓
Return success to frontend
```

### Data Distribution

| Data Type | Storage Location | Reason |
|-----------|------------------|--------|
| Task metadata (name, description, criteria) | Database | Large content, high gas cost |
| Task amount, status, milestones | Smart Contract | Core data, needs decentralization |
| Bid applications | Database | Variable quantity, high gas cost |
| Credit scores | Database | Frequent changes, high gas cost |
| Submission content (URLs) | Database | Large content, high gas cost |
| Payment records | Smart Contract | Financial data, needs transparency |
| User mnemonics | Vault | Sensitive data, needs encryption |

---

## 📊 Database Schema

### Extended Users Table
```sql
ALTER TABLE users ADD COLUMN credit_score INT DEFAULT 5000;
ALTER TABLE users ADD COLUMN tasks_completed INT DEFAULT 0;
ALTER TABLE users ADD COLUMN tasks_cancelled INT DEFAULT 0;
ALTER TABLE users ADD COLUMN eth_balance DECIMAL(20, 8) DEFAULT 0;
ALTER TABLE users ADD COLUMN xzt_balance DECIMAL(20, 8) DEFAULT 0;
```

### Tasks Table
```sql
CREATE TABLE tasks (
    task_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_task_id BIGINT UNIQUE,  -- On-chain task ID
    project_id UUID NOT NULL REFERENCES projects(project_id),
    creator_did VARCHAR(42) NOT NULL REFERENCES users(did),
    acceptor_did VARCHAR(42) REFERENCES users(did),
    
    -- Task details
    task_name VARCHAR(255) NOT NULL,
    task_description TEXT NOT NULL,
    acceptance_criteria TEXT NOT NULL,
    expected_completion_date TIMESTAMP,
    reward_amount DECIMAL(20, 8) NOT NULL,
    
    -- Visibility
    visibility VARCHAR(20) NOT NULL CHECK (visibility IN ('project', 'global')),
    
    -- Status
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP,
    design_submitted_at TIMESTAMP,
    design_approved_at TIMESTAMP,
    implementation_submitted_at TIMESTAMP,
    implementation_approved_at TIMESTAMP,
    final_submitted_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Task Bids Table
```sql
CREATE TABLE task_bids (
    bid_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(task_id),
    bidder_did VARCHAR(42) NOT NULL REFERENCES users(did),
    bid_message TEXT,
    credit_score_snapshot INT NOT NULL,  -- Credit score at time of bid
    status VARCHAR(20) NOT NULL DEFAULT 'pending' 
        CHECK (status IN ('pending', 'accepted', 'rejected')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(task_id, bidder_did)
);
```

### Task Submissions Table
```sql
CREATE TABLE task_submissions (
    submission_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(task_id),
    submission_type VARCHAR(30) NOT NULL 
        CHECK (submission_type IN ('design', 'implementation', 'final')),
    submitter_did VARCHAR(42) NOT NULL REFERENCES users(did),
    content TEXT NOT NULL,
    file_urls TEXT[],  -- Array of file URLs
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Credit History Table
```sql
CREATE TABLE credit_history (
    history_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_did VARCHAR(42) NOT NULL REFERENCES users(did),
    task_id UUID REFERENCES tasks(task_id),
    change_amount INT NOT NULL,  -- Positive for gain, negative for loss
    reason VARCHAR(50) NOT NULL,  -- 'task_completed', 'task_cancelled'
    before_score INT NOT NULL,
    after_score INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    tx_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tx_hash VARCHAR(66),  -- Blockchain transaction hash
    from_did VARCHAR(42) NOT NULL,
    to_did VARCHAR(42) NOT NULL,
    amount DECIMAL(20, 8) NOT NULL,
    tx_type VARCHAR(30) NOT NULL 
        CHECK (tx_type IN ('transfer', 'task_create', 'task_payment', 'task_refund')),
    task_id UUID REFERENCES tasks(task_id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' 
        CHECK (status IN ('pending', 'confirmed', 'failed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMP
);
```

---

## 🔌 API Endpoints

### Wallet APIs

**GET /wallet/balance**
- Description: Query user's XZT and ETH balance
- Auth: JWT token
- Response:
```json
{
  "xzt_balance": "1000.50",
  "eth_balance": "0.05",
  "did": "0x..."
}
```

**POST /wallet/transfer**
- Description: Transfer XZT to another user
- Auth: JWT token
- Request:
```json
{
  "to_address": "0x...",
  "amount": "100.00"
}
```
- Validation:
  - Amount <= 1000 XZT (individual limit)
  - System wallet: amount <= 30% of balance
- Response:
```json
{
  "tx_hash": "0x...",
  "status": "confirmed"
}
```

**GET /wallet/transactions**
- Description: Get transaction history
- Auth: JWT token
- Query params: `?page=1&limit=20&type=transfer`
- Response:
```json
{
  "transactions": [
    {
      "tx_hash": "0x...",
      "from": "0x...",
      "to": "0x...",
      "amount": "100.00",
      "type": "transfer",
      "status": "confirmed",
      "created_at": "2026-01-26T10:00:00Z"
    }
  ],
  "total": 50,
  "page": 1
}
```

---

### Task APIs

**POST /tasks**
- Description: Create a new task
- Auth: JWT token
- Request:
```json
{
  "project_id": "uuid",
  "task_name": "Build login page",
  "task_description": "...",
  "acceptance_criteria": "...",
  "expected_completion_date": "2026-02-01",
  "reward_amount": "5000.00",
  "visibility": "project"  // or "global"
}
```
- Process:
  1. Validate user has sufficient XZT balance
  2. Generate signature
  3. Call contract to create task and lock XZT
  4. Wait for transaction confirmation
  5. Save task to database
- Response:
```json
{
  "task_id": "uuid",
  "contract_task_id": 123,
  "status": "pending"
}
```

**GET /tasks**
- Description: List tasks
- Auth: JWT token
- Query params:
  - `?visibility=project&project_id=uuid` - Project tasks
  - `?visibility=global` - Global tasks
  - `?status=pending` - Filter by status
  - `?my_tasks=true` - Tasks created by me
  - `?my_bids=true` - Tasks I bid on
- Response:
```json
{
  "tasks": [
    {
      "task_id": "uuid",
      "task_name": "...",
      "reward_amount": "5000.00",
      "status": "pending",
      "creator": {
        "did": "0x...",
        "username": "alice"
      },
      "bid_count": 5,
      "created_at": "2026-01-26T10:00:00Z"
    }
  ]
}
```

**GET /tasks/:id**
- Description: Get task details
- Auth: JWT token
- Response:
```json
{
  "task_id": "uuid",
  "contract_task_id": 123,
  "task_name": "...",
  "task_description": "...",
  "acceptance_criteria": "...",
  "reward_amount": "5000.00",
  "status": "design_approved",
  "visibility": "project",
  "creator": {
    "did": "0x...",
    "username": "alice",
    "credit_score": 5000
  },
  "acceptor": {
    "did": "0x...",
    "username": "bob",
    "credit_score": 4500
  },
  "milestones": {
    "design_approved": true,
    "implementation_approved": false,
    "final_approved": false
  },
  "payments": {
    "paid_amount": "1500.00",
    "remaining_amount": "3500.00"
  },
  "submissions": [
    {
      "type": "design",
      "content": "...",
      "file_urls": ["https://..."],
      "submitted_at": "2026-01-27T10:00:00Z"
    }
  ],
  "created_at": "2026-01-26T10:00:00Z"
}
```

**POST /tasks/:id/bid**
- Description: Apply to work on a task
- Auth: JWT token
- Request:
```json
{
  "message": "I have 5 years of experience..."
}
```
- Validation:
  - Task status must be PENDING or BIDDING
  - User credit score >= 0
  - User not already bid on this task
- Response:
```json
{
  "bid_id": "uuid",
  "status": "pending"
}
```

**GET /tasks/:id/bids**
- Description: Get all bids for a task (creator only)
- Auth: JWT token
- Response:
```json
{
  "bids": [
    {
      "bid_id": "uuid",
      "bidder": {
        "did": "0x...",
        "username": "bob",
        "credit_score": 4500,
        "tasks_completed": 10,
        "tasks_cancelled": 1
      },
      "message": "...",
      "status": "pending",
      "created_at": "2026-01-26T11:00:00Z"
    }
  ]
}
```

**POST /tasks/:id/select-bidder**
- Description: Select a bidder to work on task (creator only)
- Auth: JWT token
- Request:
```json
{
  "bidder_did": "0x..."
}
```
- Process:
  1. Validate user is task creator
  2. Validate task status is BIDDING
  3. Generate signature
  4. Call contract to set acceptor
  5. Wait for confirmation
  6. Update database (accept selected bid, reject others)
- Response:
```json
{
  "status": "accepted",
  "acceptor_did": "0x..."
}
```

**POST /tasks/:id/submit**
- Description: Submit work for a milestone (executor only)
- Auth: JWT token
- Request:
```json
{
  "submission_type": "design",  // or "implementation", "final"
  "content": "Design document...",
  "file_urls": ["https://..."]
}
```
- Process:
  1. Validate user is task executor
  2. Validate current status allows submission
  3. Generate signature
  4. Call contract to update status
  5. Wait for confirmation
  6. Save submission to database
- Response:
```json
{
  "submission_id": "uuid",
  "status": "design_submitted"
}
```

**POST /tasks/:id/approve**
- Description: Approve a milestone submission (creator only)
- Auth: JWT token
- Request:
```json
{
  "milestone": "design"  // or "implementation", "final"
}
```
- Process:
  1. Validate user is task creator
  2. Validate current status
  3. Generate signature
  4. Call contract to approve milestone (triggers payment)
  5. Wait for confirmation
  6. Update database
  7. Update credit score if final approval
- Response:
```json
{
  "status": "design_approved",
  "payment": {
    "amount": "1500.00",
    "tx_hash": "0x..."
  }
}
```

**POST /tasks/:id/reject**
- Description: Reject a milestone submission (creator only)
- Auth: JWT token
- Request:
```json
{
  "milestone": "design",
  "reason": "Does not meet requirements..."
}
```
- Process:
  1. Validate user is task creator
  2. Generate signature
  3. Call contract to revert status
  4. Wait for confirmation
  5. Update database
- Response:
```json
{
  "status": "accepted",  // or previous status
  "message": "Submission rejected, executor can resubmit"
}
```

**POST /tasks/:id/cancel**
- Description: Cancel a task (creator or executor)
- Auth: JWT token
- Request:
```json
{
  "reason": "No longer needed..."
}
```
- Process:
  1. Validate user is creator or executor
  2. Calculate refund/payment based on current status
  3. Generate signature
  4. Call contract to cancel and process refund
  5. Wait for confirmation
  6. Update database
  7. Update credit score if executor cancels after milestone
- Response:
```json
{
  "status": "cancelled",
  "refund": {
    "creator_refund": "3500.00",
    "executor_payment": "1500.00",
    "tx_hash": "0x..."
  },
  "credit_penalty": -3000  // if executor cancels
}
```

---

### User APIs

**GET /users/:did/profile**
- Description: Get user profile with credit score
- Auth: JWT token
- Response:
```json
{
  "did": "0x...",
  "username": "alice",
  "credit_score": 5000,
  "tasks_completed": 10,
  "tasks_cancelled": 1,
  "xzt_balance": "1000.00",
  "member_since": "2026-01-01T00:00:00Z"
}
```

**GET /users/:did/credit-history**
- Description: Get credit score change history
- Auth: JWT token
- Response:
```json
{
  "history": [
    {
      "change_amount": -3000,
      "reason": "task_cancelled",
      "task_id": "uuid",
      "before_score": 5000,
      "after_score": 2000,
      "created_at": "2026-01-26T10:00:00Z"
    }
  ]
}
```

---

## 🚀 Implementation Plan

### Phase 1: Foundation (Week 1)

**1.1 DID Generation Migration**
- Modify `did-login-lambda/go/pkg/auth/did.go`
- Change from ed25519 to secp256k1 (Ethereum compatible)
- Use BIP44 derivation path: `m/44'/60'/0'/0/0`
- Update database schema (DID length: 66 → 42 characters)

**1.2 System Wallet Generation**
- Generate new Ethereum wallet for system
- Store mnemonic in HashiCorp Vault
- Document wallet address for initial XZT mint

**1.3 Database Schema**
- Create migration script
- Add new tables: tasks, task_bids, task_submissions, credit_history, transactions
- Extend users table with credit_score, tasks_completed, tasks_cancelled
- Create indexes for performance

**1.4 Smart Contracts**
- Write XZToken.sol (ERC20)
- Write TaskEscrow.sol with signature verification
- Write deployment scripts
- Write unit tests (Hardhat/Foundry)

---

### Phase 2: Backend APIs (Week 2)

**2.1 Wallet APIs**
- Implement balance query
- Implement transfer with limits
- Implement transaction history

**2.2 Blockchain Integration**
- Create Go package for contract interaction
- Implement signature generation
- Implement transaction sending and confirmation waiting
- Implement error handling and retry logic

**2.3 Task Management APIs**
- Implement task creation
- Implement task listing with filters
- Implement task details

**2.4 Bidding APIs**
- Implement bid submission
- Implement bid listing
- Implement bidder selection

---

### Phase 3: Task Workflow (Week 3)

**3.1 Submission APIs**
- Implement design submission
- Implement implementation submission
- Implement final submission

**3.2 Approval/Rejection APIs**
- Implement milestone approval (with payment)
- Implement milestone rejection
- Implement status rollback

**3.3 Cancellation APIs**
- Implement creator cancellation
- Implement executor cancellation
- Implement refund calculation
- Implement credit score penalty

**3.4 Credit Score System**
- Implement credit score updates
- Implement credit history tracking
- Implement bidding restrictions based on credit

---

### Phase 4: Frontend UI (Week 4)

**4.1 Wallet Page**
- Display XZT and ETH balance
- Transfer form with validation
- Transaction history table

**4.2 Task Market Page**
- Task list with filters (project/global, status)
- Task creation form
- Task card with basic info

**4.3 Task Detail Page**
- Full task information
- Bid list (for creator)
- Bid form (for executors)
- Submission forms
- Approval/rejection buttons
- Cancel button

**4.4 User Profile Page**
- Credit score display
- Task statistics
- Credit history

---

### Phase 5: Testing & Deployment (Week 5)

**5.1 Contract Deployment**
- Deploy to Sepolia testnet
- Verify contracts on Etherscan
- Mint initial 10,000 XZT to system wallet

**5.2 Backend Deployment**
- Deploy Lambda functions
- Configure API Gateway
- Set up environment variables
- Configure Vault access

**5.3 Integration Testing**
- Test complete task workflow
- Test payment flows
- Test cancellation scenarios
- Test credit score updates

**5.4 Documentation**
- API documentation
- User guide
- Admin guide
- Troubleshooting guide

---

## 🔧 Technical Specifications

### DID Generation (Updated)

**Old (ed25519)**:
```go
// 32 bytes public key → 0x + 64 hex chars = 66 chars
did := "0x" + hex.EncodeToString(ed25519PublicKey)
```

**New (secp256k1)**:
```go
// BIP44 path: m/44'/60'/0'/0/0
// Ethereum address: 20 bytes → 0x + 40 hex chars = 42 chars
privateKey := derivePrivateKey(mnemonic, "m/44'/60'/0'/0/0")
publicKey := privateKey.PublicKey()
address := crypto.PubkeyToAddress(publicKey)
did := address.Hex()  // 0x...
```

### Signature Generation

```go
func GenerateSignature(
    taskID uint64,
    action string,
    deadline int64,
    privateKey *ecdsa.PrivateKey,
) ([]byte, error) {
    // Pack data
    message := crypto.Keccak256(
        []byte(fmt.Sprintf("%d", taskID)),
        []byte(action),
        []byte(fmt.Sprintf("%d", deadline)),
    )
    
    // Sign
    signature, err := crypto.Sign(message, privateKey)
    if err != nil {
        return nil, err
    }
    
    return signature, nil
}
```

### Contract Interaction

```go
func ApproveDesign(
    taskID uint64,
    creatorDID string,
    adminPrivateKey *ecdsa.PrivateKey,
) error {
    // 1. Get user's private key from Vault
    userPrivateKey := getPrivateKeyFromVault(creatorDID)
    
    // 2. Generate signature
    deadline := time.Now().Add(5 * time.Minute).Unix()
    signature, err := GenerateSignature(taskID, "approveDesign", deadline, userPrivateKey)
    
    // 3. Build transaction
    tx := buildTransaction(
        contractAddress,
        "approveDesign",
        taskID,
        deadline,
        signature,
    )
    
    // 4. Sign transaction with admin key (for gas)
    signedTx, err := types.SignTx(tx, signer, adminPrivateKey)
    
    // 5. Send transaction
    err = client.SendTransaction(ctx, signedTx)
    
    // 6. Wait for confirmation
    receipt, err := waitForConfirmation(ctx, signedTx.Hash())
    
    // 7. Check status
    if receipt.Status != 1 {
        return errors.New("transaction failed")
    }
    
    return nil
}
```

---

## 📝 Environment Variables

```bash
# Database
DATABASE_URL=postgresql://...
DATABASE_POOLER_URL=postgresql://...

# Blockchain
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/...
CHAIN_ID=11155111
XZT_TOKEN_ADDRESS=0x...
TASK_ESCROW_ADDRESS=0x...

# Vault
VAULT_ADDR=https://vault.example.com
VAULT_TOKEN=...
VAULT_MNEMONIC_PATH=secret/data/xz-wallet/mnemonics

# Admin Wallet (for gas)
ADMIN_WALLET_ADDRESS=0x...  # Your MetaMask address
ADMIN_WALLET_PRIVATE_KEY=0x...  # Your MetaMask private key

# System Wallet
SYSTEM_WALLET_ADDRESS=0x...  # Generated during deployment

# JWT
JWT_SECRET=...

# API
API_BASE_URL=https://...
```

---

## ⚠️ Important Notes

### Gas Fee Management
- All contract calls are made by backend using admin wallet
- Admin wallet must maintain sufficient Sepolia ETH balance
- Recommended: Keep at least 0.1 ETH in admin wallet
- Monitor gas prices and adjust accordingly

### Transaction Confirmation
- Wait for at least 1 block confirmation before updating database
- Sepolia block time: ~12 seconds
- User may experience 15-30 second delay for operations
- Show "Transaction pending..." message in UI

### Error Handling
- Contract revert: Show error message to user
- Insufficient gas: Retry with higher gas price
- Network timeout: Retry up to 3 times
- Database error: Rollback transaction

### Security Considerations
- Never expose private keys in logs
- Always validate user permissions before contract calls
- Implement rate limiting on APIs
- Validate all user inputs
- Use prepared statements for SQL queries

### Scalability
- Use database connection pooling
- Cache frequently accessed data (balances, credit scores)
- Implement pagination for list APIs
- Consider using Redis for caching

---

## 📚 References

- **Ethereum**: https://ethereum.org/en/developers/docs/
- **Solidity**: https://docs.soliditylang.org/
- **OpenZeppelin**: https://docs.openzeppelin.com/contracts/
- **go-ethereum**: https://geth.ethereum.org/docs/developers/dapp-developer/native
- **BIP39**: https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
- **BIP44**: https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki

---

## ✅ Acceptance Criteria

### Functional Requirements
- ✅ Users can register and get Ethereum-compatible DID
- ✅ Users can view XZT and ETH balance
- ✅ Users can transfer XZT (with limits)
- ✅ Users can create tasks with XZT escrow
- ✅ Users can bid on tasks
- ✅ Creators can select bidders
- ✅ Executors can submit work at each milestone
- ✅ Creators can approve/reject submissions
- ✅ Automatic payment on milestone approval (30%, 80%, 100%)
- ✅ Either party can cancel with appropriate refunds
- ✅ Credit score system tracks user reputation
- ✅ Credit score affects bidding eligibility

### Non-Functional Requirements
- ✅ All payments are transparent on blockchain
- ✅ Signature verification prevents unauthorized operations
- ✅ Transaction confirmation ensures data consistency
- ✅ Private keys are securely stored in Vault
- ✅ API response time < 2 seconds (excluding blockchain confirmation)
- ✅ System can handle 100 concurrent users

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-26  
**Status**: Ready for Implementation
