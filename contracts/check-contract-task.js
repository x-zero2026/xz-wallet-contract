const { ethers } = require('ethers');
const fs = require('fs');

async function checkNextTaskId() {
    const provider = new ethers.JsonRpcProvider(process.env.SEPOLIA_RPC_URL);
    const escrowABI = JSON.parse(fs.readFileSync('./TaskEscrow.abi.json', 'utf8'));
    const escrowAddress = process.env.TASK_ESCROW_ADDRESS;
    
    const escrow = new ethers.Contract(escrowAddress, escrowABI, provider);
    const nextTaskId = await escrow.nextTaskId();
    
    console.log('Next Task ID:', nextTaskId.toString());
    console.log('Last created task ID:', (nextTaskId - 1n).toString());
    
    // Get task details
    if (nextTaskId > 0n) {
        const lastTaskId = nextTaskId - 1n;
        const task = await escrow.getTask(lastTaskId);
        console.log('\nLast Task Details:');
        console.log('Creator:', task.creator);
        console.log('Executor:', task.executor);
        console.log('Total Amount:', ethers.formatEther(task.totalAmount), 'XZT');
        console.log('Paid Amount:', ethers.formatEther(task.paidAmount), 'XZT');
        console.log('Cancelled:', task.cancelled);
    }
}

checkNextTaskId().catch(console.error);
