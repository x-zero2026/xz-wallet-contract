import { TASK_STATUS_LABELS, VISIBILITY_LABELS } from '../api';
import './TaskCard.css';

function TaskCard({ task, onClick }) {
  const getStatusClass = (status) => {
    if (status === 'pending') return 'status-pending';
    if (status === 'bidding') return 'status-bidding';
    if (status === 'accepted') return 'status-accepted';
    if (status.includes('submitted')) return 'status-submitted';
    if (status.includes('approved')) return 'status-approved';
    if (status === 'completed') return 'status-completed';
    if (status === 'cancelled') return 'status-cancelled';
    return '';
  };

  return (
    <div className="task-card" onClick={onClick}>
      <div className="task-card-header">
        <h3 className="task-card-title">{task.task_name}</h3>
        <span className={`status-badge ${getStatusClass(task.status)}`}>
          {TASK_STATUS_LABELS[task.status] || task.status}
        </span>
      </div>
      
      <p className="task-card-description">
        {task.task_description.length > 100
          ? task.task_description.substring(0, 100) + '...'
          : task.task_description}
      </p>
      
      <div className="task-card-footer">
        <div className="task-card-reward">
          <span className="reward-label">奖励:</span>
          <span className="reward-amount">{parseFloat(task.reward_amount).toFixed(2)} XZT</span>
        </div>
        <div className="task-card-visibility">
          {VISIBILITY_LABELS[task.visibility]}
        </div>
      </div>
      
      {task.executor_username && (
        <div className="task-card-executor">
          执行者: {task.executor_username}
        </div>
      )}
    </div>
  );
}

export default TaskCard;
