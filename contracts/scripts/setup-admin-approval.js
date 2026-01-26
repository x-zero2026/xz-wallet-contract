const hre = require("hardhat");
require("dotenv").config();

async function main() {
  // Get addresses from environment
  const tokenAddress = process.env.XZT_TOKEN_ADDRESS;
  const escrowAddress = process.env.TASK_ESCROW_ADDRESS;
  const adminAddress = process.env.SYSTEM_WALLET_ADDRESS;
  
  console.log("XZT Token:", tokenAddress);
  console.log("Task Escrow:", escrowAddress);
  console.log("\nAdmin Address:", adminAddress);
  
  // Get contract instances
  const XZToken = await hre.ethers.getContractAt("XZToken", tokenAddress);
  const TaskEscrow = await hre.ethers.getContractAt("TaskEscrow", escrowAddress);
  
  // Check admin's XZT balance
  const balance = await XZToken.balanceOf(adminAddress);
  console.log("Admin XZT Balance:", hre.ethers.formatEther(balance), "XZT");
  
  if (balance === 0n) {
    console.log("\n⚠️  Admin has no XZT tokens!");
    console.log("Transferring 10,000 XZT to admin...");
    
    const [deployer] = await hre.ethers.getSigners();
    const amount = hre.ethers.parseEther("10000");
    const tx = await XZToken.transfer(adminAddress, amount);
    await tx.wait();
    
    const newBalance = await XZToken.balanceOf(adminAddress);
    console.log("✅ Transferred 10,000 XZT to admin");
    console.log("New Balance:", hre.ethers.formatEther(newBalance), "XZT");
    console.log("Transaction:", tx.hash);
  }
  
  // Check current allowance
  const currentAllowance = await XZToken.allowance(adminAddress, escrowAddress);
  console.log("\nCurrent Escrow Allowance:", hre.ethers.formatEther(currentAllowance), "XZT");
  
  // Note: We cannot approve from here because we don't have admin's private key
  // The admin needs to approve the escrow contract themselves
  
  console.log("\n⚠️  IMPORTANT: Admin needs to approve escrow contract!");
  console.log("The Lambda function needs to do this approval before creating tasks.");
  console.log("\nOr run this script with admin's private key in hardhat.config.js");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
