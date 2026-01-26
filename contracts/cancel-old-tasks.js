const { ethers } = require('hardhat');

async function main() {
  const TASK_ESCROW_ADDRESS = '0x8e98B971884e14C5da6D528932bf96296311B8cb';
  
  // 获取合约实例
  const TaskEscrow = await ethers.getContractFactory('TaskEscrow');
  const escrow = TaskEscrow.attach(TASK_ESCROW_ADDRESS);
  
  // 获取管理员账户
  const [admin] = await ethers.getSigners();
  console.log(`管理员地址: ${admin.address}\n`);
  
  // 要取消的任务
  const tasksToCancel = [
    { id: 1, executorAmount: 0 }, // 任务1，未支付任何金额
    { id: 2, executorAmount: 0 }, // 任务2，未支付任何金额
  ];
  
  for (const task of tasksToCancel) {
    console.log(`========== 取消任务 #${task.id} ==========`);
    
    try {
      // 检查任务状态
      const taskInfo = await escrow.getTask(task.id);
      console.log(`创建者: ${taskInfo.creator}`);
      console.log(`总金额: ${ethers.formatUnits(taskInfo.totalAmount, 18)} XZT`);
      console.log(`已支付: ${ethers.formatUnits(taskInfo.paidAmount, 18)} XZT`);
      console.log(`已取消: ${taskInfo.cancelled}`);
      
      if (taskInfo.cancelled) {
        console.log('✅ 任务已经取消，跳过\n');
        continue;
      }
      
      // 取消任务
      console.log('正在取消任务...');
      const tx = await escrow.cancelTask(
        task.id,
        ethers.parseUnits(task.executorAmount.toString(), 18)
      );
      
      console.log(`交易哈希: ${tx.hash}`);
      console.log('等待确认...');
      
      const receipt = await tx.wait();
      console.log(`✅ 任务已取消！Gas 使用: ${receipt.gasUsed.toString()}`);
      
      // 检查创建者余额
      const XZT_TOKEN_ADDRESS = '0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8';
      const XZToken = await ethers.getContractFactory('XZToken');
      const token = XZToken.attach(XZT_TOKEN_ADDRESS);
      
      const creatorBalance = await token.balanceOf(taskInfo.creator);
      console.log(`创建者当前余额: ${ethers.formatUnits(creatorBalance, 18)} XZT\n`);
      
    } catch (error) {
      console.error(`❌ 取消任务失败: ${error.message}\n`);
    }
  }
  
  console.log('========== 完成 ==========');
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
