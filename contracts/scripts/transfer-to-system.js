const hre = require("hardhat");

async function main() {
  const tokenAddress = "0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8";
  const systemWallet = "0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA";
  const amount = hre.ethers.parseEther("10000"); // 10,000 XZT

  console.log("🔄 Transferring XZT to system wallet...\n");
  
  const [signer] = await hre.ethers.getSigners();
  console.log("From:", signer.address);
  console.log("To:", systemWallet);
  console.log("Amount:", hre.ethers.formatEther(amount), "XZT\n");

  // Get token contract
  const token = await hre.ethers.getContractAt("XZToken", tokenAddress);
  
  // Check current balances
  const fromBalance = await token.balanceOf(signer.address);
  const toBalance = await token.balanceOf(systemWallet);
  
  console.log("Current balances:");
  console.log("  Deployer:", hre.ethers.formatEther(fromBalance), "XZT");
  console.log("  System wallet:", hre.ethers.formatEther(toBalance), "XZT\n");
  
  // Transfer
  console.log("Executing transfer...");
  const tx = await token.transfer(systemWallet, amount);
  console.log("Transaction hash:", tx.hash);
  
  console.log("Waiting for confirmation...");
  const receipt = await tx.wait();
  console.log("✅ Transfer confirmed in block:", receipt.blockNumber, "\n");
  
  // Check new balances
  const newFromBalance = await token.balanceOf(signer.address);
  const newToBalance = await token.balanceOf(systemWallet);
  
  console.log("New balances:");
  console.log("  Deployer:", hre.ethers.formatEther(newFromBalance), "XZT");
  console.log("  System wallet:", hre.ethers.formatEther(newToBalance), "XZT\n");
  
  console.log("🎉 Transfer complete!");
  console.log("View on Etherscan:");
  console.log(`https://sepolia.etherscan.io/tx/${tx.hash}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ Transfer failed:", error);
    process.exit(1);
  });
