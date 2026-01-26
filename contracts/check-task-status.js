const { ethers } = require('hardhat');

async function main() {
  // 合约地址
  const TASK_ESCROW_ADDRESS = '0x8e98B971884e14C5da6D528932bf96296311B8cb';
  
  // 获取合约实例
  const TaskEscrow = await ethers.getContractFactory('TaskEscrow');
  const escrow = TaskEscrow.attach(TASK_ESCROW_ADDRESS);
  
  // 获取下一个任务ID（当前最大ID + 1）
  const nextTaskId = await escrow.nextTaskId();
  console.log(`\n下一个任务ID: ${nextTaskId}`);
  console.log(`当前已创建的任务数: ${nextTaskId - 1n}\n`);
  
  // 查询所有已创建的任务
  for (let i = 1n; i < nextTaskId; i++) {
    try {
      const task = await escrow.getTask(i);
      
      console.log(`========== 任务 #${i} ==========`);
      console.log(`创建者: ${task.creator}`);
      console.log(`执行者: ${task.executor}`);
      console.log(`总金额: ${ethers.formatUnits(task.totalAmount, 18)} XZT`);
      console.log(`已支付: ${ethers.formatUnits(task.paidAmount, 18)} XZT`);
      console.log(`剩余金额: ${ethers.formatUnits(task.totalAmount - task.paidAmount, 18)} XZT`);
      console.log(`已取消: ${task.cancelled}`);
      
      // 检查执行者余额
      if (task.executor !== '0x0000000000000000000000000000000000000000') {
        const XZT_TOKEN_ADDRESS = '0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8';
        const XZToken = await ethers.getContractFactory('XZToken');
        const token = XZToken.attach(XZT_TOKEN_ADDRESS);
        
        const executorBalance = await token.balanceOf(task.executor);
        console.log(`执行者余额: ${ethers.formatUnits(executorBalance, 18)} XZT`);
      }
      
      console.log('');
    } catch (error) {
      console.log(`任务 #${i} 不存在或查询失败: ${error.message}\n`);
    }
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
