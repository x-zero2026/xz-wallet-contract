const hre = require("hardhat");
require("dotenv").config();

async function main() {
  // Get addresses from environment
  const tokenAddress = process.env.XZT_TOKEN_ADDRESS;
  const escrowAddress = process.env.TASK_ESCROW_ADDRESS;
  const systemWalletAddress = process.env.SYSTEM_WALLET_ADDRESS;
  const systemWalletPrivateKey = process.env.SYSTEM_WALLET_PRIVATE_KEY || process.env.ADMIN_WALLET_PRIVATE_KEY;
  
  if (!systemWalletPrivateKey) {
    console.error("❌ SYSTEM_WALLET_PRIVATE_KEY not found in .env");
    process.exit(1);
  }
  
  console.log("XZT Token:", tokenAddress);
  console.log("Task Escrow:", escrowAddress);
  console.log("System Wallet:", systemWalletAddress);
  
  // Create wallet from private key
  const wallet = new hre.ethers.Wallet(systemWalletPrivateKey, hre.ethers.provider);
  console.log("\nWallet Address:", wallet.address);
  
  // Get contract instances
  const XZToken = await hre.ethers.getContractAt("XZToken", tokenAddress, wallet);
  
  // Check balance
  const balance = await XZToken.balanceOf(wallet.address);
  console.log("Current XZT Balance:", hre.ethers.formatEther(balance), "XZT");
  
  // Check current allowance
  const currentAllowance = await XZToken.allowance(wallet.address, escrowAddress);
  console.log("Current Escrow Allowance:", hre.ethers.formatEther(currentAllowance), "XZT");
  
  // Approve max amount
  console.log("\n🔄 Approving escrow contract...");
  const maxApproval = hre.ethers.MaxUint256;
  const tx = await XZToken.approve(escrowAddress, maxApproval);
  console.log("Transaction sent:", tx.hash);
  
  await tx.wait();
  console.log("✅ Approval confirmed!");
  
  // Verify new allowance
  const newAllowance = await XZToken.allowance(wallet.address, escrowAddress);
  console.log("New Escrow Allowance:", hre.ethers.formatEther(newAllowance), "XZT");
  
  console.log("\n✅ Setup complete! Admin can now create tasks.");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
