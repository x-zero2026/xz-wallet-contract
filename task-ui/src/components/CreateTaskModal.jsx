import { useState } from 'react';
import { createTask, VISIBILITY } from '../api';
import './TaskDetailModal.css'; // 使用相同的CSS样式

function CreateTaskModal({ onClose, onSuccess, selectedProject }) {
  const [formData, setFormData] = useState({
    task_name: '',
    task_description: '',
    acceptance_criteria: '',
    reward_amount: '',
    visibility: VISIBILITY.GLOBAL,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [blockchainLoading, setBlockchainLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validation
    if (!selectedProject) {
      setError('请选择项目');
      return;
    }
    if (!formData.task_name.trim()) {
      setError('请输入任务名称');
      return;
    }
    if (!formData.task_description.trim()) {
      setError('请输入任务描述');
      return;
    }
    if (!formData.acceptance_criteria.trim()) {
      setError('请输入验收标准');
      return;
    }
    if (!formData.reward_amount || parseFloat(formData.reward_amount) <= 0) {
      setError('请输入有效的奖励金额');
      return;
    }

    try {
      setLoading(true);
      setBlockchainLoading(true);
      setError(null);
      
      const response = await createTask({
        ...formData,
        project_id: selectedProject.project_id,
      });
      
      setBlockchainLoading(false);
      onSuccess(response.data.data);
      onClose();
    } catch (err) {
      setBlockchainLoading(false);
      setError(err.response?.data?.error || '创建任务失败');
      console.error('Create task error:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      {/* Blockchain Loading Overlay */}
      {blockchainLoading && (
        <div className="blockchain-loading-overlay" onClick={(e) => e.stopPropagation()}>
          <div className="blockchain-loading-content">
            <div className="blockchain-loading-spinner"></div>
            <p className="blockchain-loading-text">正在创建任务并锁定资金...</p>
          </div>
        </div>
      )}

      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">发布新任务</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>

        {error && <div className="error">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label">任务名称 *</label>
            <input
              type="text"
              name="task_name"
              className="form-input"
              value={formData.task_name}
              onChange={handleChange}
              placeholder="输入任务名称"
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label className="form-label">任务描述 *</label>
            <textarea
              name="task_description"
              className="form-textarea"
              value={formData.task_description}
              onChange={handleChange}
              placeholder="详细描述任务内容和要求"
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label className="form-label">验收标准 *</label>
            <textarea
              name="acceptance_criteria"
              className="form-textarea"
              value={formData.acceptance_criteria}
              onChange={handleChange}
              placeholder="明确的验收标准"
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label className="form-label">奖励金额 (XZT) *</label>
            <input
              type="number"
              name="reward_amount"
              className="form-input"
              value={formData.reward_amount}
              onChange={handleChange}
              placeholder="0.00"
              step="0.01"
              min="0"
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label className="form-label">可见性</label>
            <select
              name="visibility"
              className="form-select"
              value={formData.visibility}
              onChange={handleChange}
              disabled={loading}
            >
              <option value={VISIBILITY.GLOBAL}>全局可见</option>
              <option value={VISIBILITY.PROJECT}>项目内可见</option>
            </select>
          </div>

          <div className="modal-footer">
            <button
              type="button"
              className="btn btn-secondary"
              onClick={onClose}
              disabled={loading}
            >
              取消
            </button>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading}
            >
              {loading ? '发布中...' : '发布任务'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default CreateTaskModal;
