# XZ Wallet Documentation

Complete documentation for the XZ Wallet & Task Escrow System.

## 📖 Documentation Index

### Getting Started

- **[QUICKSTART.md](./QUICKSTART.md)** - 5-minute quick start guide
- **[DEPLOYMENT-GUIDE.md](./DEPLOYMENT-GUIDE.md)** - Complete deployment instructions
- **[INTEGRATION.md](./INTEGRATION.md)** - Integration with DID login system

### Project Information

- **[PROGRESS.md](./PROGRESS.md)** - Development progress and milestones
- **[CONFIGURATION.md](./CONFIGURATION.md)** - System configuration summary
- **[SUMMARY.md](./SUMMARY.md)** - Task Center implementation summary

### Requirements

- **[REQUIREMENTS.md](./REQUIREMENTS.md)** - Complete requirements specification
- **[REQUIREMENTS-MVP.md](./REQUIREMENTS-MVP.md)** - MVP requirements and scope

### Deployment History

- **[DEPLOYMENT-SUCCESS.md](./DEPLOYMENT-SUCCESS.md)** - Successful deployment records

## 🏗️ Architecture

### Smart Contracts (Sepolia Testnet)

- **XZToken**: ERC20 token contract
  - Address: `0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8`
  - [View on Etherscan](https://sepolia.etherscan.io/address/0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8)

- **TaskEscrow**: Task management and escrow contract
  - Address: `0x8e98B971884e14C5da6D528932bf96296311B8cb`
  - [View on Etherscan](https://sepolia.etherscan.io/address/0x8e98B971884e14C5da6D528932bf96296311B8cb)

### Backend (AWS Lambda)

- **API Gateway**: https://yms07x0sn0.execute-api.us-east-1.amazonaws.com/prod/
- **Lambda Functions**:
  - GetBalanceFunction - Query XZT balance
  - CreateTaskFunction - Create new task
  - ListTasksFunction - List tasks with filters
  - GetTaskFunction - Get task details
  - BidTaskFunction - Bid on task
  - SelectBidderFunction - Select executor
  - ApproveWorkFunction - Approve work and pay

### Frontend (React)

- **Task Center**: Task management UI
- **Development**: http://localhost:5173
- **Test Login**: http://localhost:5173/test-login.html

### Database (PostgreSQL/Supabase)

- Users, tasks, bids, submissions, credit history
- Schema: [../database/schema.sql](../database/schema.sql)

## 📚 Component Documentation

### Smart Contracts
- [Contracts README](../contracts/README.md)
- [XZToken.sol](../contracts/contracts/XZToken.sol)
- [TaskEscrow.sol](../contracts/contracts/TaskEscrow.sol)

### Backend
- [Lambda README](../lambda/README.md)
- [API Documentation](../lambda/README.md#api-endpoints)

### Frontend
- [Task UI README](../task-ui/README.md)
- [Integration Guide](./INTEGRATION.md)
- [Quick Start](./QUICKSTART.md)

### Scripts
- [Scripts README](../scripts/README.md)
- [register-app.sh](../scripts/register-app.sh)

## 🔄 Task Workflow

```
1. Creator publishes task
   ↓ XZT locked in contract
2. Executors bid on task
   ↓ Credit score displayed
3. Creator selects executor
   ↓ Task status: accepted
4. Executor submits design
   ↓ Creator approves → 30% paid
5. Executor submits implementation
   ↓ Creator approves → 50% paid (80% total)
6. Executor submits final work
   ↓ Creator approves → 20% paid (100% total)
7. Task completed
```

## 🏆 Credit Score System

| Event | Points Change |
|-------|--------------|
| Initial score | 5000 |
| Complete task | +100 |
| Cancel after design approved | -3000 |
| Cancel after implementation approved | -8000 |

**Rules**:
- Credit score < 0: Cannot bid on tasks
- Credit score affects bidding competitiveness

## 🔐 Security

- User authentication via JWT tokens
- Smart contract enforced payments
- Milestone-based escrow
- Credit score reputation system
- Admin wallet for gas fees

## 📊 Current Status

| Component | Status | Progress |
|-----------|--------|----------|
| Smart Contracts | ✅ Deployed | 100% |
| Database | ✅ Migrated | 100% |
| Backend (Lambda) | ✅ Deployed | 100% |
| Frontend (Task UI) | ✅ Complete | 100% |
| Integration | ✅ Complete | 100% |
| **Overall** | **✅ Complete** | **100%** |

## 🚀 Next Steps

1. Test complete task workflow end-to-end
2. Deploy Task Center to production (AWS Amplify)
3. Add more Lambda functions (submit-work, reject-work, cancel-task)
4. Implement real-time notifications
5. Add task search and filtering
6. Mobile responsive optimization

## 📞 Support

For questions or issues:
1. Check the relevant documentation above
2. Review [PROGRESS.md](./PROGRESS.md) for current status
3. See [DEPLOYMENT-GUIDE.md](./DEPLOYMENT-GUIDE.md) for deployment help

---

**Last Updated**: 2026-01-26  
**Version**: 1.0
