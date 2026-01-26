# XZ Wallet Deployment Guide

Complete guide for deploying XZ Wallet smart contracts and backend to production.

## 📋 Prerequisites Checklist

- [ ] Node.js 18+ installed
- [ ] Go 1.21+ installed
- [ ] MetaMask wallet with Sepolia ETH (at least 0.1 ETH)
- [ ] Infura or Alchemy account (for Sepolia RPC)
- [ ] Etherscan API key (for contract verification)
- [ ] AWS account with Lambda access
- [ ] Supabase PostgreSQL database
- [ ] HashiCorp Vault instance

## 🚀 Step-by-Step Deployment

### Step 1: Deploy Smart Contracts

#### 1.1 Install Dependencies

```bash
cd xz-wallet-contract/contracts
npm install
```

#### 1.2 Configure Environment

```bash
cp .env.example .env
```

Edit `.env`:
```env
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
ADMIN_WALLET_PRIVATE_KEY=0x...  # Your MetaMask private key
ETHERSCAN_API_KEY=YOUR_ETHERSCAN_API_KEY
```

#### 1.3 Compile Contracts

```bash
npm run compile
```

Expected output:
```
Compiled 2 Solidity files successfully
```

#### 1.4 Deploy to Sepolia

```bash
npm run deploy:sepolia
```

Expected output:
```
🚀 Starting deployment to Sepolia...
📝 Deploying contracts with account: 0x...
💰 Account balance: 0.1 ETH

📦 Deploying XZToken...
✅ XZToken deployed to: 0x...
   Initial supply: 10000.0 XZT

📦 Deploying TaskEscrow...
✅ TaskEscrow deployed to: 0x...

💾 Deployment info saved to: deployment.json
```

#### 1.5 Verify Contracts on Etherscan

```bash
# Verify XZToken
npx hardhat verify --network sepolia <XZT_TOKEN_ADDRESS>

# Verify TaskEscrow
npx hardhat verify --network sepolia <TASK_ESCROW_ADDRESS> <XZT_TOKEN_ADDRESS>
```

#### 1.6 Save Contract Addresses

Update `contracts/.env`:
```env
XZT_TOKEN_ADDRESS=0x...
TASK_ESCROW_ADDRESS=0x...
```

---

### Step 2: Setup Database

#### 2.1 Connect to Supabase

```bash
psql "postgresql://postgres.rbpsksuuvtzmathnmyxn:iPass4xz2026!@aws-1-ap-south-1.pooler.supabase.com:6543/postgres"
```

#### 2.2 Run Migration

```bash
cd ../database
psql "postgresql://..." -f schema.sql
```

Expected output:
```
ALTER TABLE
CREATE INDEX
CREATE TABLE
...
```

#### 2.3 Verify Tables

```sql
\dt
```

Should show:
- users (extended with credit_score, xzt_balance)
- tasks
- task_bids
- task_submissions
- credit_history
- xzt_transactions

---

### Step 3: Generate Go Bindings

#### 3.1 Install abigen

```bash
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

#### 3.2 Generate Bindings

```bash
cd ../lambda
make generate-bindings
```

Expected output:
```
Generating Go bindings from contract ABIs...
✅ Go bindings generated successfully
```

This creates:
- `pkg/blockchain/contracts/xztoken.go`
- `pkg/blockchain/contracts/taskescrow.go`

---

### Step 4: Build Lambda Functions

#### 4.1 Install Go Dependencies

```bash
go mod download
```

#### 4.2 Build All Functions

```bash
make build
```

Expected output:
```
Building Lambda functions...
Building cmd/create-task...
Building cmd/pay-milestone...
...
✅ All Lambda functions built
```

---

### Step 5: Deploy to AWS Lambda

#### 5.1 Create SAM Template

Create `lambda/template.yaml`:

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31

Globals:
  Function:
    Timeout: 30
    MemorySize: 128
    Runtime: provided.al2
    Architectures:
      - arm64
    Environment:
      Variables:
        # Database
        DATABASE_URL: !Ref DatabaseURL
        DB_PASSWORD: !Ref DBPassword
        
        # Blockchain
        SEPOLIA_RPC_URL: !Ref SepoliaRPCURL
        CHAIN_ID: "11155111"
        XZT_TOKEN_ADDRESS: !Ref XZTTokenAddress
        TASK_ESCROW_ADDRESS: !Ref TaskEscrowAddress
        
        # Admin Wallet
        ADMIN_WALLET_ADDRESS: !Ref AdminWalletAddress
        ADMIN_WALLET_PRIVATE_KEY: !Ref AdminWalletPrivateKey
        
        # Vault
        VAULT_ADDR: !Ref VaultAddr
        VAULT_TOKEN: !Ref VaultToken
        
        # JWT
        JWT_SECRET: !Ref JWTSecret

Parameters:
  DatabaseURL:
    Type: String
  DBPassword:
    Type: String
    NoEcho: true
  SepoliaRPCURL:
    Type: String
  XZTTokenAddress:
    Type: String
  TaskEscrowAddress:
    Type: String
  AdminWalletAddress:
    Type: String
  AdminWalletPrivateKey:
    Type: String
    NoEcho: true
  VaultAddr:
    Type: String
  VaultToken:
    Type: String
    NoEcho: true
  JWTSecret:
    Type: String
    NoEcho: true

Resources:
  CreateTaskFunction:
    Type: AWS::Serverless::Function
    Properties:
      CodeUri: cmd/create-task/
      Handler: bootstrap
      Events:
        CreateTask:
          Type: Api
          Properties:
            Path: /tasks
            Method: post

  # Add more functions...
```

#### 5.2 Deploy with SAM

```bash
sam build
sam deploy --guided
```

Follow prompts to configure:
- Stack name: `xz-wallet-backend`
- AWS Region: `us-east-1`
- Confirm changes: `Y`
- Allow SAM CLI IAM role creation: `Y`
- Save arguments to config: `Y`

---

### Step 6: Configure API Gateway

#### 6.1 Enable CORS

In AWS Console:
1. Go to API Gateway
2. Select your API
3. Enable CORS for all endpoints
4. Deploy API

#### 6.2 Get API Endpoint

```bash
aws apigateway get-rest-apis --query 'items[?name==`xz-wallet-api`].id' --output text
```

Save the endpoint URL:
```
https://xxxxx.execute-api.us-east-1.amazonaws.com/prod
```

---

### Step 7: Test Deployment

#### 7.1 Test Contract Interaction

```bash
# Check XZT balance
curl https://your-api.com/wallet/balance \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 7.2 Create Test Task

```bash
curl -X POST https://your-api.com/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "uuid",
    "task_name": "Test Task",
    "task_description": "Test description",
    "acceptance_criteria": "Test criteria",
    "reward_amount": "100.00",
    "visibility": "global"
  }'
```

#### 7.3 Verify on Blockchain

Check transaction on Sepolia Etherscan:
```
https://sepolia.etherscan.io/address/<TASK_ESCROW_ADDRESS>
```

---

## 🔧 Configuration Reference

### Environment Variables

#### Smart Contracts
```env
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_KEY
ADMIN_WALLET_PRIVATE_KEY=0x...
ETHERSCAN_API_KEY=...
XZT_TOKEN_ADDRESS=0x...
TASK_ESCROW_ADDRESS=0x...
```

#### Lambda Functions
```env
# Database
DATABASE_URL=postgresql://...
DB_PASSWORD=...

# Blockchain
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_KEY
CHAIN_ID=11155111
XZT_TOKEN_ADDRESS=0x...
TASK_ESCROW_ADDRESS=0x...

# Admin Wallet
ADMIN_WALLET_ADDRESS=0x...
ADMIN_WALLET_PRIVATE_KEY=0x...

# Vault
VAULT_ADDR=http://your-vault-server:8200
VAULT_TOKEN=your_vault_token_here

# JWT
JWT_SECRET=your_jwt_secret_here
```

---

## 🐛 Troubleshooting

### Contract Deployment Issues

**Error: Insufficient funds**
- Get Sepolia ETH from faucets:
  - https://sepoliafaucet.com/
  - https://www.alchemy.com/faucets/ethereum-sepolia

**Error: Nonce too low**
- Wait a few seconds and retry
- Or reset MetaMask account

### Lambda Deployment Issues

**Error: Invalid architecture**
- Ensure building with `GOARCH=arm64`
- Check Lambda configuration uses `arm64`

**Error: Function timeout**
- Increase timeout in SAM template
- Check Sepolia RPC is responding

### Database Issues

**Error: Connection refused**
- Check Supabase pooler URL
- Verify password is correct
- Check IP whitelist

---

## ✅ Post-Deployment Checklist

- [ ] Contracts deployed and verified on Etherscan
- [ ] 10,000 XZT minted to system wallet
- [ ] Database schema migrated
- [ ] Lambda functions deployed
- [ ] API Gateway configured with CORS
- [ ] Environment variables set
- [ ] Test task created successfully
- [ ] Blockchain transactions confirmed
- [ ] API endpoints responding
- [ ] Frontend can connect to API

---

## 📝 Next Steps

1. **Transfer XZT to users** for testing
2. **Monitor gas usage** and optimize if needed
3. **Set up CloudWatch alarms** for Lambda errors
4. **Configure backup** for Vault
5. **Document API endpoints** for frontend team
6. **Create admin dashboard** for monitoring

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-26  
**Status**: Ready for Deployment
