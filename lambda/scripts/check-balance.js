const { ethers } = require('ethers');

// 配置
const RPC_URL = 'https://eth-sepolia.g.alchemy.com/v2/GA5ibaTuz122ssPqQhWL7';
const XZT_TOKEN_ADDRESS = '0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8';

// XZT Token ABI (只需要 balanceOf)
const XZT_ABI = [
  'function balanceOf(address account) view returns (uint256)',
  'function decimals() view returns (uint8)'
];

async function checkBalance(address) {
  const provider = new ethers.JsonRpcProvider(RPC_URL);
  const token = new ethers.Contract(XZT_TOKEN_ADDRESS, XZT_ABI, provider);
  
  const balance = await token.balanceOf(address);
  const decimals = await token.decimals();
  
  const formattedBalance = ethers.formatUnits(balance, decimals);
  
  console.log(`Address: ${address}`);
  console.log(`Balance: ${formattedBalance} XZT`);
  console.log(`Raw Balance: ${balance.toString()}`);
}

// 从命令行参数获取地址
const address = process.argv[2];
if (!address) {
  console.error('Usage: node check-balance.js <address>');
  process.exit(1);
}

checkBalance(address).catch(console.error);
