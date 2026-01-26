# XZ Wallet & Task Escrow System

A decentralized task management and payment system built on Ethereum Sepolia testnet with milestone-based escrow and credit score mechanism.

## 🎯 Overview

- **Blockchain**: Sepolia Testnet
- **Token**: XZT (ERC20) - 1 XZT = 1 e-CNY
- **Backend**: Go + AWS Lambda
- **Database**: PostgreSQL (Supabase)
- **Key Management**: HashiCorp Vault

## 📋 Key Features

- **Wallet Management**: Ethereum-compatible DID with托管钱包
- **Token System**: XZT token with transfer limits
- **Task Escrow**: Milestone-based payments (30%, 80%, 100%)
- **Bidding System**: Competitive bidding with credit score display
- **Credit Score**: Reputation system (5000 initial, penalties for cancellation)
- **Decentralized**: Smart contract enforced payments and state transitions

## 📚 Documentation

### Quick Links
- **[PROGRESS.md](./docs/PROGRESS.md)** - Development progress and status
- **[CONFIGURATION.md](./docs/CONFIGURATION.md)** - System configuration summary
- **[DEPLOYMENT-GUIDE.md](./docs/DEPLOYMENT-GUIDE.md)** - Deployment instructions

### Requirements
- **[REQUIREMENTS.md](./docs/REQUIREMENTS.md)** - Complete requirements specification
- **[REQUIREMENTS-MVP.md](./docs/REQUIREMENTS-MVP.md)** - MVP requirements

### Task Center (Frontend)
- **[INTEGRATION.md](./docs/INTEGRATION.md)** - Integration with DID login system
- **[QUICKSTART.md](./docs/QUICKSTART.md)** - Quick start guide
- **[SUMMARY.md](./docs/SUMMARY.md)** - Implementation summary

### Components
- **[Contracts README](./contracts/README.md)** - Smart contract documentation
- **[Lambda README](./lambda/README.md)** - Backend API documentation
- **[Task UI README](./task-ui/README.md)** - Frontend documentation
- **[Scripts README](./scripts/README.md)** - Utility scripts

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL (Supabase)
- AWS CLI & SAM CLI
- MetaMask wallet

### 1. Deploy Smart Contracts

```bash
cd contracts
npm install
cp .env.example .env
# Edit .env with your configuration

# Compile contracts
npm run compile

# Deploy to Sepolia
npm run deploy
```

### 2. Setup Database

```bash
# Run schema migration on Supabase
psql -h your-supabase-host -U postgres -d postgres -f database/schema.sql
```

### 3. Deploy Backend (Lambda)

```bash
cd lambda
cp .env.deploy.example .env.deploy
# Edit .env.deploy with your configuration

# Build and deploy
sam build
sam deploy --parameter-overrides $(cat .env.deploy | tr '\n' ' ')
```

### 4. Run Task Center Frontend

```bash
cd task-ui
npm install
cp .env.example .env
# Edit .env with API URLs

# Start development server
npm run dev
```

Visit http://localhost:5173/test-login.html to test login.

### 5. Register App in DID Login System

```bash
cd scripts
export JWT_TOKEN="your-jwt-token"
./register-app.sh
```

For detailed instructions, see [DEPLOYMENT-GUIDE.md](./docs/DEPLOYMENT-GUIDE.md)

## 📊 Project Structure

```
xz-wallet-contract/
├── contracts/          # Smart contracts (Solidity)
│   ├── contracts/
│   │   ├── XZToken.sol
│   │   └── TaskEscrow.sol
│   ├── scripts/
│   └── test/
├── lambda/             # Backend APIs (Go + AWS Lambda)
│   ├── cmd/           # Lambda functions
│   ├── pkg/           # Shared packages
│   └── template.yaml  # SAM template
├── task-ui/           # Task Center Frontend (React)
│   ├── src/
│   │   ├── components/
│   │   ├── api/
│   │   └── utils/
│   └── public/
├── database/          # Database schemas
│   └── schema.sql
├── scripts/           # Utility scripts
│   └── register-app.sh
└── docs/             # Documentation
    ├── REQUIREMENTS.md
    ├── PROGRESS.md
    ├── CONFIGURATION.md
    ├── DEPLOYMENT-GUIDE.md
    ├── INTEGRATION.md
    ├── QUICKSTART.md
    └── SUMMARY.md
```

## 🔄 Task Workflow

```
1. Creator publishes task → XZT locked in contract
2. Executors bid on task
3. Creator selects executor
4. Executor submits design → Creator approves → 30% paid
5. Executor submits implementation → Creator approves → 50% paid (80% total)
6. Executor submits final work → Creator approves → 20% paid (100% total)
```

## 🏆 Credit Score System

- Initial: 5000 points
- Cancel after design approved: -3000 points
- Cancel after implementation approved: -8000 points
- Complete task: +100 points
- Credit < 0: Cannot bid on tasks

## 🔐 Security

- User mnemonics stored in HashiCorp Vault
- Signature verification for all contract operations
- Backend代签 with user authorization
- Admin wallet (MetaMask) pays gas fees

## 📝 License

MIT

## 👥 Team

X-Zero Platform Development Team

---

**Status**: ✅ Backend Deployed & Tested | 🔄 Frontend Complete  
**Version**: 1.0  
**Last Updated**: 2026-01-26

## 🔗 Deployed Resources

- **XZToken Contract**: [0x6b1f...98F8](https://sepolia.etherscan.io/address/0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8)
- **TaskEscrow Contract**: [0x8e98...8cb](https://sepolia.etherscan.io/address/0x8e98B971884e14C5da6D528932bf96296311B8cb)
- **API Gateway**: https://yms07x0sn0.execute-api.us-east-1.amazonaws.com/prod/
- **Task Center**: http://localhost:5173 (development)
