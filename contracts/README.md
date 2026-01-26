# XZ Wallet Smart Contracts

Smart contracts for XZ Wallet task escrow system on Sepolia testnet.

## 📋 Contracts

### XZToken.sol
- Standard ERC20 token
- Symbol: XZT
- Initial supply: 10,000 XZT
- Decimals: 18

### TaskEscrow.sol
- Escrow contract for task management
- Milestone-based payments (30%, 80%, 100%)
- Task cancellation with refunds
- Admin-controlled (MVP version)

## 🚀 Setup

### Prerequisites
- Node.js 18+
- npm or yarn
- MetaMask wallet with Sepolia ETH

### Installation

```bash
cd xz-wallet-contract/contracts
npm install
```

### Configuration

1. Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

2. Fill in your credentials:
```env
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
ADMIN_WALLET_PRIVATE_KEY=0x...
ETHERSCAN_API_KEY=YOUR_ETHERSCAN_API_KEY
```

## 📦 Compilation

```bash
npm run compile
```

This will:
- Compile contracts
- Generate ABIs in `artifacts/` directory
- Generate TypeScript types

## 🚀 Deployment

### Deploy to Sepolia

```bash
npm run deploy:sepolia
```

This will:
1. Deploy XZToken contract
2. Deploy TaskEscrow contract
3. Mint 10,000 XZT to deployer
4. Save deployment info to `deployment.json`

### Verify on Etherscan

```bash
# Verify XZToken
npx hardhat verify --network sepolia <XZT_TOKEN_ADDRESS>

# Verify TaskEscrow
npx hardhat verify --network sepolia <TASK_ESCROW_ADDRESS> <XZT_TOKEN_ADDRESS>
```

## 🧪 Testing

```bash
npm test
```

## 📝 Contract Addresses (Sepolia)

After deployment, update these in your `.env`:

```env
XZT_TOKEN_ADDRESS=0x...
TASK_ESCROW_ADDRESS=0x...
```

## 🔐 Security Notes

- Never commit `.env` file
- Keep private keys secure
- Admin wallet pays all gas fees
- Only admin can call escrow functions

## 📚 Next Steps

After deployment:

1. **Generate Go bindings**:
   ```bash
   cd ../lambda
   make generate-bindings
   ```

2. **Update Lambda environment variables** with contract addresses

3. **Transfer XZT to system wallet** if needed

4. **Test contract interactions** from backend

## 🛠️ Development

### Local Testing

```bash
# Start local Hardhat node
npx hardhat node

# Deploy to local network
npx hardhat run scripts/deploy.js --network localhost
```

### Contract Interaction

```javascript
const { ethers } = require("hardhat");

// Get contract instance
const token = await ethers.getContractAt("XZToken", tokenAddress);
const escrow = await ethers.getContractAt("TaskEscrow", escrowAddress);

// Check balance
const balance = await token.balanceOf(address);
console.log("Balance:", ethers.formatEther(balance), "XZT");

// Create task
const tx = await escrow.createTask(creator, executor, amount);
await tx.wait();
```

## 📖 API Reference

### XZToken

```solidity
function balanceOf(address account) external view returns (uint256)
function transfer(address to, uint256 amount) external returns (bool)
function approve(address spender, uint256 amount) external returns (bool)
function transferFrom(address from, address to, uint256 amount) external returns (bool)
function mint(address to, uint256 amount) external onlyOwner
function burn(uint256 amount) external
```

### TaskEscrow

```solidity
function createTask(address creator, address executor, uint256 amount) external onlyOwner returns (uint256)
function setExecutor(uint256 taskId, address executor) external onlyOwner
function payMilestone(uint256 taskId, uint256 amount) external onlyOwner
function cancelTask(uint256 taskId, uint256 executorAmount) external onlyOwner
function getTask(uint256 taskId) external view returns (...)
function getRemainingAmount(uint256 taskId) external view returns (uint256)
```

## 🐛 Troubleshooting

### Deployment fails with "insufficient funds"
- Ensure your admin wallet has enough Sepolia ETH
- Get free Sepolia ETH from faucets:
  - https://sepoliafaucet.com/
  - https://www.alchemy.com/faucets/ethereum-sepolia

### Contract verification fails
- Check Etherscan API key is correct
- Wait a few minutes after deployment
- Ensure constructor arguments match

### Transaction reverts
- Check admin wallet is set as owner
- Ensure sufficient XZT balance
- Verify task exists and is not cancelled

## 📄 License

MIT
