const hre = require("hardhat");
require("dotenv").config();

async function main() {
  const [deployer] = await hre.ethers.getSigners();
  
  // Get addresses from environment
  const tokenAddress = process.env.XZT_TOKEN_ADDRESS;
  const escrowAddress = process.env.TASK_ESCROW_ADDRESS;
  const adminAddress = process.env.SYSTEM_WALLET_ADDRESS;
  const adminPrivateKey = process.env.ADMIN_WALLET_PRIVATE_KEY;
  
  console.log("Deployer:", deployer.address);
  console.log("XZT Token:", tokenAddress);
  console.log("Task Escrow:", escrowAddress);
  console.log("Admin Address:", adminAddress);
  
  // Get token contract
  const XZToken = await hre.ethers.getContractAt("XZToken", tokenAddress);
  
  // Check admin's current balance
  const adminBalance = await XZToken.balanceOf(adminAddress);
  console.log("\nAdmin XZT Balance:", hre.ethers.formatEther(adminBalance), "XZT");
  
  // Transfer 10,000 XZT to admin if needed
  if (adminBalance < hre.ethers.parseEther("10000")) {
    console.log("\n🔄 Transferring 10,000 XZT to admin...");
    const amount = hre.ethers.parseEther("10000");
    const tx = await XZToken.transfer(adminAddress, amount);
    console.log("Transaction sent:", tx.hash);
    await tx.wait();
    console.log("✅ Transfer confirmed!");
    
    const newBalance = await XZToken.balanceOf(adminAddress);
    console.log("New Admin Balance:", hre.ethers.formatEther(newBalance), "XZT");
  }
  
  // Now approve escrow using admin's private key
  console.log("\n🔄 Approving escrow contract with admin wallet...");
  const adminWallet = new hre.ethers.Wallet(adminPrivateKey, hre.ethers.provider);
  const XZTokenAsAdmin = await hre.ethers.getContractAt("XZToken", tokenAddress, adminWallet);
  
  const maxApproval = hre.ethers.MaxUint256;
  const approveTx = await XZTokenAsAdmin.approve(escrowAddress, maxApproval);
  console.log("Approval transaction sent:", approveTx.hash);
  await approveTx.wait();
  console.log("✅ Approval confirmed!");
  
  // Verify allowance
  const allowance = await XZToken.allowance(adminAddress, escrowAddress);
  console.log("\nFinal Escrow Allowance:", hre.ethers.formatEther(allowance), "XZT");
  
  console.log("\n✅ Setup complete! Admin can now create tasks.");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
