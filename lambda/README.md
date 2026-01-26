# XZ Wallet Lambda Functions

Go-based AWS Lambda functions for XZ Wallet backend.

## 📁 Project Structure

```
lambda/
├── cmd/                    # Lambda function handlers
│   ├── get-balance/       # Get XZT balance
│   ├── transfer-xzt/      # Transfer XZT
│   ├── create-task/       # Create task and lock XZT
│   ├── list-tasks/        # List tasks
│   ├── get-task/          # Get task details
│   ├── bid-task/          # Bid on task
│   ├── select-bidder/     # Select bidder
│   ├── submit-work/       # Submit work
│   ├── approve-work/      # Approve work and pay milestone
│   ├── reject-work/       # Reject work
│   └── cancel-task/       # Cancel task with refund
├── pkg/                   # Shared packages
│   ├── blockchain/        # Blockchain client and contract interaction
│   │   ├── client.go     # Ethereum client wrapper
│   │   ├── escrow.go     # TaskEscrow operations
│   │   └── contracts/    # Generated contract bindings
│   ├── models/           # Data models
│   │   └── task.go       # Task-related models
│   ├── db/               # Database connection
│   │   └── postgres.go   # PostgreSQL pool
│   ├── auth/             # Authentication
│   │   └── jwt.go        # JWT validation
│   └── response/         # API responses
│       └── response.go   # Response helpers
├── go.mod                # Go dependencies
├── go.sum                # Dependency checksums
├── Makefile              # Build commands
└── README.md             # This file
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- `abigen` tool (from go-ethereum)
- AWS CLI configured
- SAM CLI (for deployment)

### Setup

1. **Generate contract bindings**:
```bash
make generate-bindings
```

2. **Install dependencies**:
```bash
go mod download
```

3. **Build all functions**:
```bash
make build
```

## 📦 Lambda Functions

### Wallet Functions

#### GET /wallet/balance
Get XZT balance for authenticated user.

**Headers**: `Authorization: Bearer <JWT>`

**Response**:
```json
{
  "success": true,
  "data": {
    "did": "0x...",
    "eth_address": "0x...",
    "xzt_balance": "1000.50000000",
    "username": "alice"
  }
}
```

#### POST /wallet/transfer
Transfer XZT to another user (admin only for MVP).

**Headers**: `Authorization: Bearer <JWT>`

**Request**:
```json
{
  "to_address": "0x...",
  "amount": "100.00"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "tx_hash": "0x...",
    "from": "0x...",
    "to": "0x...",
    "amount": "100.00"
  }
}
```

### Task Functions

#### POST /tasks
Create a new task and lock XZT.

**Headers**: `Authorization: Bearer <JWT>`

**Request**:
```json
{
  "project_id": "uuid",
  "task_name": "Build login page",
  "task_description": "...",
  "acceptance_criteria": "...",
  "reward_amount": "5000.00",
  "visibility": "project"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "task_id": "uuid",
    "contract_task_id": 123,
    "tx_hash": "0x...",
    "status": "pending"
  }
}
```

#### GET /tasks
List tasks with filters.

**Query Parameters**:
- `visibility`: `project` | `global`
- `project_id`: UUID (required if visibility=project)
- `status`: Task status
- `creator_did`: Filter by creator
- `executor_did`: Filter by executor

**Response**:
```json
{
  "success": true,
  "data": {
    "tasks": [
      {
        "task_id": "uuid",
        "task_name": "...",
        "reward_amount": "5000.00",
        "status": "pending",
        "creator": {...},
        "bid_count": 5
      }
    ]
  }
}
```

#### GET /tasks/:id
Get task details.

**Response**:
```json
{
  "success": true,
  "data": {
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
}
```

#### POST /tasks/:id/bid
Bid on a task.

**Headers**: `Authorization: Bearer <JWT>`

**Request**:
```json
{
  "message": "I have 5 years experience..."
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "bid_id": "uuid",
    "status": "pending"
  }
}
```

#### POST /tasks/:id/select-bidder
Select a bidder (creator only).

**Headers**: `Authorization: Bearer <JWT>`

**Request**:
```json
{
  "bidder_did": "0x..."
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "task_id": "uuid",
    "executor_did": "0x...",
    "tx_hash": "0x...",
    "status": "accepted"
  }
}
```

#### POST /tasks/:id/submit
Submit work (executor only).

**Headers**: `Authorization: Bearer <JWT>`

**Request**:
```json
{
  "submission_type": "design",
  "content": "...",
  "file_urls": ["https://..."]
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "submission_id": "uuid",
    "status": "design_submitted"
  }
}
```

#### POST /tasks/:id/approve
Approve work and pay milestone (creator only).

**Headers**: `Authorization: Bearer <JWT>`

**Request**:
```json
{
  "milestone": "design"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "design_approved",
    "payment": {
      "amount": "1500.00",
      "tx_hash": "0x..."
    }
  }
}
```

#### POST /tasks/:id/reject
Reject work (creator only).

**Headers**: `Authorization: Bearer <JWT>`

**Request**:
```json
{
  "milestone": "design",
  "reason": "Does not meet requirements"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "accepted"
  }
}
```

#### POST /tasks/:id/cancel
Cancel task with refund.

**Headers**: `Authorization: Bearer <JWT>`

**Request**:
```json
{
  "reason": "Requirements changed"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "cancelled",
    "refund": {
      "creator_refund": "3500.00",
      "executor_payment": "1500.00",
      "tx_hash": "0x..."
    },
    "credit_penalty": -3000
  }
}
```

## 🔧 Environment Variables

All Lambda functions require these environment variables:

```env
# Database
DATABASE_URL=postgresql://...
DB_PASSWORD=...

# Blockchain
SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/...
CHAIN_ID=11155111
XZT_TOKEN_ADDRESS=0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8
TASK_ESCROW_ADDRESS=0x8e98B971884e14C5da6D528932bf96296311B8cb

# Admin Wallet
ADMIN_WALLET_ADDRESS=0xd62F159A744df11332F8F1C73C827aed8Ca9378D
ADMIN_WALLET_PRIVATE_KEY=0x...

# JWT
JWT_SECRET=your_jwt_secret_here
JWT_EXPIRY=168h
```

## 🏗️ Build & Deploy

### Build All Functions

```bash
make build
```

This compiles all Lambda functions for ARM64 architecture.

### Deploy with SAM

```bash
sam build
sam deploy --guided
```

### Deploy Individual Function

```bash
cd cmd/get-balance
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
zip function.zip bootstrap
aws lambda update-function-code \
  --function-name xz-wallet-GetBalanceFunction \
  --zip-file fileb://function.zip \
  --region us-east-1
```

## 🧪 Testing

### Local Testing with SAM

```bash
sam local start-api
```

### Test Individual Function

```bash
sam local invoke GetBalanceFunction -e events/get-balance.json
```

### Create Test Event

```json
{
  "httpMethod": "GET",
  "headers": {
    "Authorization": "Bearer eyJhbGc..."
  },
  "body": null
}
```

## 📊 Monitoring

### View Logs

```bash
sam logs -n GetBalanceFunction --stack-name xz-wallet-backend --tail
```

### CloudWatch Metrics

- Invocations
- Duration
- Errors
- Throttles

## 🐛 Troubleshooting

### Common Issues

**Error: "DATABASE_URL not set"**
- Ensure environment variables are configured in Lambda

**Error: "failed to connect to Ethereum client"**
- Check SEPOLIA_RPC_URL is correct
- Verify Alchemy/Infura API key is valid

**Error: "invalid token"**
- Check JWT_SECRET matches between services
- Verify token hasn't expired

**Error: "transaction failed"**
- Check admin wallet has enough ETH for gas
- Verify contract addresses are correct

## 📝 Development

### Adding a New Function

1. Create directory: `cmd/my-function/`
2. Create `main.go` with handler
3. Update `Makefile` if needed
4. Add to SAM template
5. Build and deploy

### Code Style

- Use `gofmt` for formatting
- Follow Go best practices
- Add error handling
- Include logging

## 📚 References

- [AWS Lambda Go](https://github.com/aws/aws-lambda-go)
- [go-ethereum](https://github.com/ethereum/go-ethereum)
- [pgx](https://github.com/jackc/pgx)

---

**Status**: In Development  
**Version**: 1.0.0  
**Last Updated**: 2026-01-26
