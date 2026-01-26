# Admin Setup Guide for Task Creation

## Problem
When admin user tries to create a task, they get error: "Failed to create task on blockchain: transaction failed"

## Root Cause
The smart contract's `createTask` function requires:
1. The creator (admin) must have sufficient XZT tokens
2. The creator must have approved the TaskEscrow contract to spend their tokens

## Solution

### Step 1: Check Admin's XZT Balance
```bash
# Using the get-balance API
curl -X GET "https://yms07x0sn0.execute-api.us-east-1.amazonaws.com/prod/wallet/balance?address=0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Step 2: Transfer XZT to Admin (if needed)
If admin has no XZT, transfer from the deployer wallet:

```javascript
// Using Hardhat console or script
const XZToken = await ethers.getContractAt("XZToken", "0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8");
const tx = await XZToken.transfer("0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA", ethers.parseEther("10000"));
await tx.wait();
```

### Step 3: Approve Escrow Contract
The admin wallet needs to approve the TaskEscrow contract:

```javascript
// Get admin's private key from vault
const adminPrivateKey = "ADMIN_PRIVATE_KEY_FROM_VAULT";
const wallet = new ethers.Wallet(adminPrivateKey, ethers.provider);

// Approve escrow
const XZToken = await ethers.getContractAt("XZToken", "0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8", wallet);
const escrowAddress = "0x8e98B971884e14C5da6D528932bf96296311B8cb";
const maxApproval = ethers.MaxUint256;

const tx = await XZToken.approve(escrowAddress, maxApproval);
await tx.wait();
console.log("✅ Approved!");
```

## Contract Addresses (Sepolia)
- XZT Token: `0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8`
- Task Escrow: `0x8e98B971884e14C5da6D528932bf96296311B8cb`
- Admin Address: `0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA`

## Alternative: Use System Transfer Script
Run the provided script (requires deployer access):

```bash
cd xz-wallet-contract/contracts
node scripts/fund-and-approve-admin.js
```

## Verification
After setup, verify the allowance:

```javascript
const XZToken = await ethers.getContractAt("XZToken", "0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8");
const allowance = await XZToken.allowance(
  "0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA",  // admin
  "0x8e98B971884e14C5da6D528932bf96296311B8cb"   // escrow
);
console.log("Allowance:", ethers.formatEther(allowance), "XZT");
```

## Notes
- This is a one-time setup per user
- In production, users should approve the escrow contract themselves through the UI
- The approval is unlimited (MaxUint256) so it doesn't need to be repeated
