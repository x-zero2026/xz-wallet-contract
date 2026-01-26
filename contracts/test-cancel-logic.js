const { ethers } = require('hardhat');

async function main() {
  const TASK_ESCROW_ADDRESS = '0x8e98B971884e14C5da6D528932bf96296311B8cb';
  const XZT_TOKEN_ADDRESS = '0x6b1f7209E08Bd8B9ec44DDb4Edd9B4AA6acd98F8';
  
  const TaskEscrow = await ethers.getContractFactory('TaskEscrow');
  const escrow = TaskEscrow.attach(TASK_ESCROW_ADDRESS);
  
  const XZToken = await ethers.getContractFactory('XZToken');
  const token = XZToken.attach(XZT_TOKEN_ADDRESS);
  
  const taskId = 3;
  
  console.log('========== 任务 #3 取消前状态 ==========\n');
  
  // 获取任务信息
  const task = await escrow.getTask(taskId);
  console.log(`创建者: ${task.creator}`);
  console.log(`执行者: ${task.executor}`);
  console.log(`总金额: ${ethers.formatUnits(task.totalAmount, 18)} XZT`);
  console.log(`已支付: ${ethers.formatUnits(task.paidAmount, 18)} XZT`);
  console.log(`剩余金额: ${ethers.formatUnits(task.totalAmount - task.paidAmount, 18)} XZT`);
  console.log(`已取消: ${task.cancelled}\n`);
  
  // 获取当前余额
  const creatorBalanceBefore = await token.balanceOf(task.creator);
  const executorBalanceBefore = await token.balanceOf(task.executor);
  
  console.log('========== 当前余额 ==========\n');
  console.log(`创建者余额: ${ethers.formatUnits(creatorBalanceBefore, 18)} XZT`);
  console.log(`执行者余额: ${ethers.formatUnits(executorBalanceBefore, 18)} XZT\n`);
  
  console.log('========== 取消后预期 ==========\n');
  const executorAmount = task.paidAmount;  // 30 XZT
  const remaining = task.totalAmount - task.paidAmount;  // 70 XZT
  const creatorRefund = remaining - executorAmount;  // 70 - 30 = 40 XZT? 不对！
  
  // 正确的逻辑：
  // executorAmount 是要额外支付给执行者的金额（在已支付基础上）
  // 但我们的 Lambda 传的是 paidAmount（已经支付过的）
  // 所以实际上：
  // - 执行者不会再收到额外的钱（因为已经收到过 30 XZT）
  // - 创建者会收到全部剩余的 70 XZT
  
  console.log(`执行者额外收到: 0 XZT (已经收到过 ${ethers.formatUnits(task.paidAmount, 18)} XZT)`);
  console.log(`创建者退款: ${ethers.formatUnits(remaining, 18)} XZT`);
  console.log(`\n预期最终余额:`);
  console.log(`创建者: ${ethers.formatUnits(creatorBalanceBefore + remaining, 18)} XZT`);
  console.log(`执行者: ${ethers.formatUnits(executorBalanceBefore, 18)} XZT (不变)\n`);
  
  console.log('⚠️  注意：Lambda 传递的 executorAmount 参数有问题！');
  console.log('应该传 0，而不是 paidAmount (30 XZT)');
  console.log('因为执行者已经收到过 30 XZT 了，不应该再支付。\n');
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
