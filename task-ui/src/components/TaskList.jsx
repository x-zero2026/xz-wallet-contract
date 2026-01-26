import { useState, useEffect } from 'react';
import { listTasks, TASK_STATUS_LABELS } from '../api';
import { getUserInfo } from '../utils/auth';
import TaskCard from './TaskCard';
import './TaskList.css';

function TaskList({ filter, onTaskClick }) {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadTasks();
  }, [filter]);

  const loadTasks = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // 从 filter 中提取 exclude_creator 和 exclude_bidded 标记
      const { exclude_creator, exclude_bidded, ...apiFilter } = filter;
      
      const response = await listTasks(apiFilter);
      console.log('API response:', response);
      
      // API returns: { success: true, data: { tasks: [...], total: 5 } }
      const responseData = response?.data?.data || response?.data || {};
      console.log('Response data:', responseData);
      
      // Extract tasks array from response
      let taskArray = responseData.tasks || [];
      console.log('Task array:', taskArray);
      
      const userInfo = getUserInfo();
      
      // 如果需要排除自己发布的任务
      if (exclude_creator && userInfo?.did) {
        taskArray = taskArray.filter(task => task.creator_did !== userInfo.did);
        console.log('Filtered tasks (excluded own):', taskArray);
      }
      
      // 如果需要排除已投标的任务
      if (exclude_bidded && userInfo?.did) {
        // 获取用户已投标的任务列表
        try {
          const biddedResponse = await listTasks({ bidder_did: userInfo.did });
          const biddedTasks = biddedResponse?.data?.data?.tasks || [];
          const biddedTaskIds = new Set(biddedTasks.map(t => t.task_id));
          
          taskArray = taskArray.filter(task => !biddedTaskIds.has(task.task_id));
          console.log('Filtered tasks (excluded bidded):', taskArray);
        } catch (err) {
          console.error('Failed to get bidded tasks:', err);
          // 如果获取失败，继续显示所有任务
        }
      }
      
      // Ensure taskArray is an array
      if (Array.isArray(taskArray)) {
        setTasks(taskArray);
      } else {
        console.error('Tasks is not an array:', taskArray);
        setTasks([]);
      }
    } catch (err) {
      console.error('Load tasks error:', err);
      setError(err?.error || err?.message || '加载任务失败');
      setTasks([]); // Set empty array on error
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="loading">加载中...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  if (tasks.length === 0) {
    return (
      <div className="empty-state">
        <div className="empty-state-icon">📋</div>
        <p>暂无任务</p>
      </div>
    );
  }

  return (
    <div className="task-list">
      {tasks.map((task) => (
        <TaskCard key={task.task_id} task={task} onClick={() => onTaskClick(task)} />
      ))}
    </div>
  );
}

export default TaskList;
