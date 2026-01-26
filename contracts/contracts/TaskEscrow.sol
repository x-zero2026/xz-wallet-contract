// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title TaskEscrow
 * @dev Escrow contract for XZ Wallet task management system
 * 
 * Features:
 * - Lock XZT tokens for tasks
 * - Milestone-based payments (30%, 80%, 100%)
 * - Task cancellation with refunds
 * - Admin-controlled (MVP version)
 */
contract TaskEscrow is Ownable, ReentrancyGuard {
    
    IERC20 public immutable token;
    
    struct Task {
        address creator;      // Task creator
        address executor;     // Task executor (can be 0x0 initially)
        uint256 totalAmount;  // Total XZT locked
        uint256 paidAmount;   // Amount already paid to executor
        bool cancelled;       // Whether task is cancelled
    }
    
    // taskId => Task
    mapping(uint256 => Task) public tasks;
    
    // Next task ID (auto-increment)
    uint256 public nextTaskId;
    
    // Events
    event TaskCreated(
        uint256 indexed taskId,
        address indexed creator,
        address indexed executor,
        uint256 amount
    );
    
    event ExecutorSet(
        uint256 indexed taskId,
        address indexed executor
    );
    
    event MilestonePaid(
        uint256 indexed taskId,
        address indexed executor,
        uint256 amount,
        uint256 totalPaid
    );
    
    event TaskCancelled(
        uint256 indexed taskId,
        uint256 executorAmount,
        uint256 creatorRefund
    );
    
    /**
     * @dev Constructor
     * @param _token Address of XZT token contract
     */
    constructor(address _token) Ownable(msg.sender) {
        require(_token != address(0), "Invalid token address");
        token = IERC20(_token);
    }
    
    /**
     * @dev Create a new task and lock XZT
     * @param creator Address of task creator
     * @param executor Address of executor (can be 0x0 if not selected yet)
     * @param amount Amount of XZT to lock (in wei)
     * @return taskId The ID of created task
     */
    function createTask(
        address creator,
        address executor,
        uint256 amount
    ) external onlyOwner nonReentrant returns (uint256) {
        require(creator != address(0), "Invalid creator");
        require(amount > 0, "Amount must be positive");
        
        // Transfer XZT from creator to this contract
        require(
            token.transferFrom(creator, address(this), amount),
            "Transfer failed"
        );
        
        uint256 taskId = nextTaskId++;
        
        tasks[taskId] = Task({
            creator: creator,
            executor: executor,
            totalAmount: amount,
            paidAmount: 0,
            cancelled: false
        });
        
        emit TaskCreated(taskId, creator, executor, amount);
        
        return taskId;
    }
    
    /**
     * @dev Set or update executor for a task
     * @param taskId ID of the task
     * @param executor Address of new executor
     */
    function setExecutor(
        uint256 taskId,
        address executor
    ) external onlyOwner {
        require(taskId < nextTaskId, "Task does not exist");
        require(executor != address(0), "Invalid executor");
        
        Task storage task = tasks[taskId];
        require(!task.cancelled, "Task is cancelled");
        
        task.executor = executor;
        
        emit ExecutorSet(taskId, executor);
    }
    
    /**
     * @dev Pay milestone to executor
     * @param taskId ID of the task
     * @param amount Amount to pay (in wei)
     */
    function payMilestone(
        uint256 taskId,
        uint256 amount
    ) external onlyOwner nonReentrant {
        require(taskId < nextTaskId, "Task does not exist");
        
        Task storage task = tasks[taskId];
        require(!task.cancelled, "Task is cancelled");
        require(task.executor != address(0), "No executor set");
        require(amount > 0, "Amount must be positive");
        require(
            task.paidAmount + amount <= task.totalAmount,
            "Exceeds total amount"
        );
        
        task.paidAmount += amount;
        
        require(
            token.transfer(task.executor, amount),
            "Transfer failed"
        );
        
        emit MilestonePaid(taskId, task.executor, amount, task.paidAmount);
    }
    
    /**
     * @dev Cancel task with refund distribution
     * @param taskId ID of the task
     * @param executorAmount Amount to pay executor (in wei)
     */
    function cancelTask(
        uint256 taskId,
        uint256 executorAmount
    ) external onlyOwner nonReentrant {
        require(taskId < nextTaskId, "Task does not exist");
        
        Task storage task = tasks[taskId];
        require(!task.cancelled, "Already cancelled");
        
        uint256 remaining = task.totalAmount - task.paidAmount;
        require(executorAmount <= remaining, "Exceeds remaining amount");
        
        task.cancelled = true;
        task.paidAmount = task.totalAmount; // Mark as fully paid
        
        // Pay executor if applicable
        if (executorAmount > 0 && task.executor != address(0)) {
            require(
                token.transfer(task.executor, executorAmount),
                "Executor payment failed"
            );
        }
        
        // Refund creator
        uint256 creatorRefund = remaining - executorAmount;
        if (creatorRefund > 0) {
            require(
                token.transfer(task.creator, creatorRefund),
                "Creator refund failed"
            );
        }
        
        emit TaskCancelled(taskId, executorAmount, creatorRefund);
    }
    
    /**
     * @dev Get task details
     * @param taskId ID of the task
     */
    function getTask(uint256 taskId) external view returns (
        address creator,
        address executor,
        uint256 totalAmount,
        uint256 paidAmount,
        bool cancelled
    ) {
        require(taskId < nextTaskId, "Task does not exist");
        Task memory task = tasks[taskId];
        return (
            task.creator,
            task.executor,
            task.totalAmount,
            task.paidAmount,
            task.cancelled
        );
    }
    
    /**
     * @dev Get remaining amount for a task
     * @param taskId ID of the task
     */
    function getRemainingAmount(uint256 taskId) external view returns (uint256) {
        require(taskId < nextTaskId, "Task does not exist");
        Task memory task = tasks[taskId];
        return task.totalAmount - task.paidAmount;
    }
    
    /**
     * @dev Emergency withdraw (owner only, for stuck tokens)
     * @param to Address to send tokens
     * @param amount Amount to withdraw
     */
    function emergencyWithdraw(
        address to,
        uint256 amount
    ) external onlyOwner {
        require(to != address(0), "Invalid address");
        require(token.transfer(to, amount), "Transfer failed");
    }
}
