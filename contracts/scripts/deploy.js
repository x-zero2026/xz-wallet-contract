const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🚀 Starting deployment to Sepolia...\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("📝 Deploying contracts with account:", deployer.address);
  
  const balance = await hre.ethers.provider.getBalance(deployer.address);
  console.log("💰 Account balance:", hre.ethers.formatEther(balance), "ETH\n");

  // Deploy XZToken
  console.log("📦 Deploying XZToken...");
  const XZToken = await hre.ethers.getContractFactory("XZToken");
  const token = await XZToken.deploy();
  await token.waitForDeployment();
  const tokenAddress = await token.getAddress();
  console.log("✅ XZToken deployed to:", tokenAddress);
  
  // Check initial supply
  const totalSupply = await token.totalSupply();
  console.log("   Initial supply:", hre.ethers.formatEther(totalSupply), "XZT");
  const deployerBalance = await token.balanceOf(deployer.address);
  console.log("   Deployer balance:", hre.ethers.formatEther(deployerBalance), "XZT\n");

  // Deploy TaskEscrow
  console.log("📦 Deploying TaskEscrow...");
  const TaskEscrow = await hre.ethers.getContractFactory("TaskEscrow");
  const escrow = await TaskEscrow.deploy(tokenAddress);
  await escrow.waitForDeployment();
  const escrowAddress = await escrow.getAddress();
  console.log("✅ TaskEscrow deployed to:", escrowAddress);
  console.log("   Token address:", await escrow.token());
  console.log("   Owner:", await escrow.owner(), "\n");

  // Save deployment info
  const deploymentInfo = {
    network: "sepolia",
    chainId: 11155111,
    deployer: deployer.address,
    deployedAt: new Date().toISOString(),
    contracts: {
      XZToken: {
        address: tokenAddress,
        initialSupply: hre.ethers.formatEther(totalSupply)
      },
      TaskEscrow: {
        address: escrowAddress,
        tokenAddress: tokenAddress
      }
    }
  };

  const deploymentPath = path.join(__dirname, "../deployment.json");
  fs.writeFileSync(deploymentPath, JSON.stringify(deploymentInfo, null, 2));
  console.log("💾 Deployment info saved to:", deploymentPath, "\n");

  // Print summary
  console.log("=" .repeat(60));
  console.log("📋 DEPLOYMENT SUMMARY");
  console.log("=".repeat(60));
  console.log("Network:        Sepolia Testnet");
  console.log("Deployer:      ", deployer.address);
  console.log("XZToken:       ", tokenAddress);
  console.log("TaskEscrow:    ", escrowAddress);
  console.log("Initial Supply:", hre.ethers.formatEther(totalSupply), "XZT");
  console.log("=".repeat(60), "\n");

  console.log("📝 Next steps:");
  console.log("1. Verify contracts on Etherscan:");
  console.log(`   npx hardhat verify --network sepolia ${tokenAddress}`);
  console.log(`   npx hardhat verify --network sepolia ${escrowAddress} ${tokenAddress}`);
  console.log("\n2. Update .env file with contract addresses");
  console.log("\n3. Generate Go bindings:");
  console.log("   cd ../lambda");
  console.log("   make generate-bindings");
  console.log("\n4. Transfer XZT to system wallet if needed");
  console.log("\n✨ Deployment complete!");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ Deployment failed:", error);
    process.exit(1);
  });
