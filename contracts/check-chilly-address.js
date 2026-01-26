const { ethers } = require('hardhat');

async function main() {
  const XZT_TOKEN_ADDRESS = '0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8';
  
  // chilly 可能的地址
  const addresses = [
    '0xED9AF503187F0F8FB1ff6BceFedc4C6bA215ebE9', // 任务 #3 的执行者
    // 如果有其他可能的地址，可以添加在这里
  ];
  
  const XZToken = await ethers.getContractFactory('XZToken');
  const token = XZToken.attach(XZT_TOKEN_ADDRESS);
  
  console.log('\n检查各地址的 XZT 余额:\n');
  
  for (const address of addresses) {
    const balance = await token.balanceOf(address);
    console.log(`地址: ${address}`);
    console.log(`余额: ${ethers.formatUnits(balance, 18)} XZT\n`);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
