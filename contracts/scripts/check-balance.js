const hre = require("hardhat");

async function main() {
  const [signer] = await hre.ethers.getSigners();
  console.log("Address:", signer.address);
  
  const balance = await hre.ethers.provider.getBalance(signer.address);
  console.log("Balance:", hre.ethers.formatEther(balance), "ETH");
  
  if (balance === 0n) {
    console.log("\n⚠️  Warning: Account has no ETH!");
    console.log("Get Sepolia ETH from:");
    console.log("- https://sepoliafaucet.com/");
    console.log("- https://www.alchemy.com/faucets/ethereum-sepolia");
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
