import { useState, useEffect } from 'react';
import { getTask, bidTask, selectBidder, approveWork, cancelTask, submitWork, TASK_STATUS_LABELS } from '../api';
import { getUserInfo } from '../utils/auth';
import './TaskDetailModal.css';

function TaskDetailModal({ taskId, onClose, onUpdate, onSwitchTab }) {
  const [task, setTask] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [actionLoading, setActionLoading] = useState(false);
  const [bidMessage, setBidMessage] = useState('');
  const [showBidForm, setShowBidForm] = useState(false);
  const [showSubmitForm, setShowSubmitForm] = useState(false);
  const [submitContent, setSubmitContent] = useState('');
  const [submitType, setSubmitType] = useState('');
  const [selectingBidder, setSelectingBidder] = useState(false);
  const [approvingWork, setApprovingWork] = useState(false);
  const [blockchainLoading, setBlockchainLoading] = useState(false);
  const [blockchainMessage, setBlockchainMessage] = useState('');
  
  const userInfo = getUserInfo();
  const isCreator = task && task.creator_did === userInfo?.did;
  const isExecutor = task && task.executor_did === userInfo?.did;

  useEffect(() => {
    loadTask();
  }, [taskId]);

  const loadTask = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await getTask(taskId);
      const taskData = response.data.data;
      // API returns {task, creator, executor, submissions, bids}
      // We need to merge task with additional info
      setTask({
        ...taskData.task,
        creator: taskData.creator,
        executor: taskData.executor,
        submissions: taskData.submissions || [],
        bids: taskData.bids || []
      });
    } catch (err) {
      setError(err.response?.data?.error || '加载任务详情失败');
      console.error('Load task error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleBid = async () => {
    try {
      setActionLoading(true);
      setBlockchainLoading(true);
      setBlockchainMessage('正在提交投标...');
      setError(null);
      await bidTask(taskId, { message: bidMessage });
      setBlockchainLoading(false);
      alert('投标成功！');
      setShowBidForm(false);
      setBidMessage('');
      onUpdate();
      onClose();
      // 切换到"我的任务"标签
      if (onSwitchTab) {
        onSwitchTab('my-tasks');
      }
    } catch (err) {
      setBlockchainLoading(false);
      setError(err.response?.data?.error || '投标失败');
      console.error('Bid error:', err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleSelectBidder = async (bidderDid) => {
    if (!confirm('确定选择此投标者吗？')) return;
    
    try {
      setSelectingBidder(true);
      setBlockchainLoading(true);
      setBlockchainMessage('正在选择执行者并锁定资金...');
      setError(null);
      await selectBidder(taskId, { bidder_did: bidderDid });
      setBlockchainLoading(false);
      alert('已选择执行者！');
      onUpdate();
      onClose(); // 关闭弹窗
    } catch (err) {
      setBlockchainLoading(false);
      setError(err.response?.data?.error || '选择执行者失败');
      console.error('Select bidder error:', err);
    } finally {
      setSelectingBidder(false);
    }
  };

  const handleApprove = async (submissionType, approve) => {
    const action = approve ? '批准' : '拒绝';
    if (!confirm(`确定${action}此提交吗？`)) return;
    
    try {
      setApprovingWork(true);
      if (approve) {
        setBlockchainLoading(true);
        // 根据不同阶段显示不同的提示文字
        let message = '正在支付里程碑奖励...';
        if (submissionType === 'design') {
          message = '正在支付设计奖励...';
        } else if (submissionType === 'implementation') {
          message = '正在支付基础成果奖励...';
        } else if (submissionType === 'final') {
          message = '正在支付最终奖励...';
        }
        setBlockchainMessage(message);
      }
      setError(null);
      await approveWork(taskId, {
        milestone: submissionType,
        approve: approve,
      });
      setBlockchainLoading(false);
      alert(`${action}成功！`);
      onUpdate();
      onClose(); // 关闭弹窗
    } catch (err) {
      setBlockchainLoading(false);
      setError(err.response?.data?.error || `${action}失败`);
      console.error('Approve error:', err);
    } finally {
      setApprovingWork(false);
    }
  };

  const handleCancelTask = async () => {
    if (!confirm('确定要取消此任务吗？取消后将无法恢复。')) return;
    
    try {
      setActionLoading(true);
      setBlockchainLoading(true);
      setBlockchainMessage('正在取消任务并退款...');
      setError(null);
      await cancelTask(taskId);
      setBlockchainLoading(false);
      alert('任务已取消');
      onUpdate();
      onClose();
    } catch (err) {
      setBlockchainLoading(false);
      setError(err.response?.data?.error || '取消任务失败');
      console.error('Cancel task error:', err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleQuitTask = async () => {
    // Determine penalty based on current status
    let penaltyMessage = '';
    if (task.status === 'design_approved' || task.status === 'implementation_submitted') {
      penaltyMessage = '设计已完成，退出将扣除100信用分。';
    } else if (task.status === 'implementation_approved' || task.status === 'final_submitted') {
      penaltyMessage = '基础成果已完成，退出将扣除200信用分。';
    }
    
    const confirmMessage = penaltyMessage 
      ? `${penaltyMessage}\n\n确定要退出此任务吗？退出后将无法恢复。`
      : '确定要退出此任务吗？退出后将无法恢复。';
    
    if (!confirm(confirmMessage)) return;
    
    try {
      setActionLoading(true);
      setBlockchainLoading(true);
      setBlockchainMessage('正在退出任务并处理退款...');
      setError(null);
      await cancelTask(taskId);
      setBlockchainLoading(false);
      alert('已退出任务');
      onUpdate();
      onClose();
    } catch (err) {
      setBlockchainLoading(false);
      setError(err.response?.data?.error || '退出任务失败');
      console.error('Quit task error:', err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleSubmitWork = async () => {
    if (!submitContent.trim()) {
      setError('请填写提交内容');
      return;
    }

    try {
      setActionLoading(true);
      setBlockchainLoading(true);
      setBlockchainMessage('正在提交工作成果...');
      setError(null);
      await submitWork(taskId, {
        submission_type: submitType,
        content: submitContent,
      });
      setBlockchainLoading(false);
      alert('提交成功！等待创建者审批');
      setShowSubmitForm(false);
      setSubmitContent('');
      setSubmitType('');
      onUpdate();
      onClose(); // 关闭弹窗
    } catch (err) {
      setBlockchainLoading(false);
      setError(err.response?.data?.error || '提交失败');
      console.error('Submit work error:', err);
    } finally {
      setActionLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="modal-overlay" onClick={onClose}>
        <div className="modal" onClick={(e) => e.stopPropagation()}>
          <div className="loading">加载中...</div>
        </div>
      </div>
    );
  }

  if (!task) {
    return null;
  }

  const getStatusClass = (status) => {
    if (!status) return '';
    if (status === 'pending') return 'status-pending';
    if (status === 'bidding') return 'status-bidding';
    if (status === 'accepted') return 'status-accepted';
    if (status.includes('submitted')) return 'status-submitted';
    if (status.includes('approved')) return 'status-approved';
    if (status === 'completed') return 'status-completed';
    if (status === 'cancelled') return 'status-cancelled';
    return '';
  };

  const canBid = task.status === 'bidding' && !isCreator && !isExecutor;
  const canSelectBidder = isCreator && task.status === 'bidding';
  const canCancelTask = isCreator && (task.status !== 'completed' && task.status !== 'cancelled');
  const canQuitTask = isExecutor && (task.status !== 'completed' && task.status !== 'cancelled');
  const canApproveDesign = isCreator && task.status === 'design_submitted';
  const canApproveImplementation = isCreator && task.status === 'implementation_submitted';
  const canApproveFinal = isCreator && task.status === 'final_submitted';
  
  // Executor submit actions
  const canSubmitDesign = isExecutor && task.status === 'accepted';
  const canSubmitImplementation = isExecutor && task.status === 'design_approved';
  const canSubmitFinal = isExecutor && task.status === 'implementation_approved';

  return (
    <div className="modal-overlay" onClick={onClose}>
      {/* Blockchain Loading Overlay */}
      {blockchainLoading && (
        <div className="blockchain-loading-overlay" onClick={(e) => e.stopPropagation()}>
          <div className="blockchain-loading-content">
            <div className="blockchain-loading-spinner"></div>
            <p className="blockchain-loading-text">{blockchainMessage}</p>
          </div>
        </div>
      )}

      <div className="modal modal-large" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">{task.task_name}</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>

        {error && <div className="error">{error}</div>}

        <div className="task-detail">
          <div className="detail-section">
            <div className="detail-row">
              <span className="detail-label">状态:</span>
              <span className={`status-badge ${getStatusClass(task.status)}`}>
                {TASK_STATUS_LABELS[task.status]}
              </span>
            </div>
            <div className="detail-row">
              <span className="detail-label">奖励:</span>
              <span className="detail-value reward">{parseFloat(task.reward_amount).toFixed(2)} XZT</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">已支付:</span>
              <span className="detail-value">{parseFloat(task.paid_amount).toFixed(2)} XZT</span>
            </div>
          </div>

          <div className="detail-section">
            <h3 className="section-title">任务描述</h3>
            <p className="detail-text">{task.task_description}</p>
          </div>

          <div className="detail-section">
            <h3 className="section-title">验收标准</h3>
            <p className="detail-text">{task.acceptance_criteria}</p>
          </div>

          {task.executor && (
            <div className="detail-section">
              <h3 className="section-title">执行者</h3>
              <div className="detail-text">
                <div>
                  <strong>{task.executor.username}</strong>
                  <span style={{ marginLeft: '1rem', color: '#1890ff' }}>
                    信用分: {task.executor.credit_score}
                  </span>
                </div>
              </div>
            </div>
          )}

          {/* Submissions History (for creator and executor to view) */}
          {(isCreator || isExecutor) && task.submissions && task.submissions.length > 0 && (
            <div className="detail-section">
              <h3 className="section-title">提交历史</h3>
              <div className="submissions-list">
                {task.submissions.map((submission) => (
                  <div key={submission.submission_id} className="submission-item">
                    <div className="submission-header">
                      <span className="submission-type">
                        {submission.submission_type === 'design' && '设计方案'}
                        {submission.submission_type === 'implementation' && '基础成果'}
                        {submission.submission_type === 'final' && '最终成果'}
                      </span>
                      <span className={`submission-status status-${submission.status}`}>
                        {submission.status === 'pending' && '待审批'}
                        {submission.status === 'approved' && '已批准'}
                        {submission.status === 'rejected' && '已拒绝'}
                      </span>
                    </div>
                    <div className="submission-content">
                      {submission.content}
                    </div>
                    {submission.file_urls && submission.file_urls.length > 0 && (
                      <div className="submission-files">
                        <strong>附件:</strong>
                        {submission.file_urls.map((url, idx) => (
                          <a key={idx} href={url} target="_blank" rel="noopener noreferrer" className="file-link">
                            文件 {idx + 1}
                          </a>
                        ))}
                      </div>
                    )}
                    <div className="submission-time">
                      提交时间: {new Date(submission.submitted_at).toLocaleString('zh-CN')}
                    </div>
                    {submission.rejection_reason && (
                      <div className="submission-rejection">
                        拒绝原因: {submission.rejection_reason}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Bidders List (for creator when status is bidding) */}
          {canSelectBidder && task.bids && task.bids.length > 0 && (
            <div className="detail-section">
              <h3 className="section-title">投标列表 ({task.bids.length})</h3>
              <div className="bids-list">
                {task.bids.map((bid) => (
                  <div key={bid.bid_id} className="bid-item">
                    <div className="bid-info">
                      <div className="bid-bidder">
                        <strong>{bid.bidder_username}</strong>
                        <span className="bid-credit">信用分: {bid.bidder_credit_score}</span>
                      </div>
                      <div className="bid-message">
                        {bid.bid_message || '(未填写申请说明)'}
                      </div>
                      <div className="bid-time">
                        投标时间: {new Date(bid.created_at).toLocaleString('zh-CN')}
                      </div>
                    </div>
                    <button
                      className="btn btn-primary btn-sm"
                      onClick={() => handleSelectBidder(bid.bidder_did)}
                      disabled={selectingBidder}
                    >
                      {selectingBidder ? '选择中...' : '选择'}
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Bid Actions */}
          {canBid && (
            <div className="detail-section">
              {!showBidForm ? (
                <button
                  className="btn btn-primary"
                  onClick={() => setShowBidForm(true)}
                  disabled={actionLoading}
                >
                  申请任务
                </button>
              ) : (
                <div className="bid-form">
                  <textarea
                    className="form-textarea"
                    placeholder="说明你的优势和计划（可选）"
                    value={bidMessage}
                    onChange={(e) => setBidMessage(e.target.value)}
                    disabled={actionLoading}
                  />
                  <div className="bid-form-actions">
                    <button
                      className="btn btn-secondary"
                      onClick={() => {
                        setShowBidForm(false);
                        setBidMessage('');
                      }}
                      disabled={actionLoading}
                    >
                      取消
                    </button>
                    <button
                      className="btn btn-primary"
                      onClick={handleBid}
                      disabled={actionLoading}
                    >
                      {actionLoading ? '提交中...' : '提交申请'}
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Approval Actions */}
          {canApproveDesign && (
            <div className="detail-section">
              <h3 className="section-title">设计方案审批</h3>
              <div className="approval-actions">
                <button
                  className="btn btn-success"
                  onClick={() => handleApprove('design', true)}
                  disabled={approvingWork}
                >
                  批准设计
                </button>
                <button
                  className="btn btn-danger"
                  onClick={() => handleApprove('design', false)}
                  disabled={approvingWork}
                >
                  拒绝设计
                </button>
              </div>
            </div>
          )}

          {canApproveImplementation && (
            <div className="detail-section">
              <h3 className="section-title">基础成果审批</h3>
              <div className="approval-actions">
                <button
                  className="btn btn-success"
                  onClick={() => handleApprove('implementation', true)}
                  disabled={approvingWork}
                >
                  批准基础成果
                </button>
                <button
                  className="btn btn-danger"
                  onClick={() => handleApprove('implementation', false)}
                  disabled={approvingWork}
                >
                  拒绝基础成果
                </button>
              </div>
            </div>
          )}

          {canApproveFinal && (
            <div className="detail-section">
              <h3 className="section-title">最终成果审批</h3>
              <div className="approval-actions">
                <button
                  className="btn btn-success"
                  onClick={() => handleApprove('final', true)}
                  disabled={approvingWork}
                >
                  批准并完成任务
                </button>
                <button
                  className="btn btn-danger"
                  onClick={() => handleApprove('final', false)}
                  disabled={approvingWork}
                >
                  拒绝最终成果
                </button>
              </div>
            </div>
          )}

          {/* Executor Submit Actions */}
          {(canSubmitDesign || canSubmitImplementation || canSubmitFinal) && (
            <div className="detail-section">
              {!showSubmitForm ? (
                <button
                  className="btn btn-primary"
                  onClick={() => {
                    setShowSubmitForm(true);
                    if (canSubmitDesign) setSubmitType('design');
                    else if (canSubmitImplementation) setSubmitType('implementation');
                    else if (canSubmitFinal) setSubmitType('final');
                  }}
                  disabled={actionLoading}
                >
                  {canSubmitDesign && '提交设计方案'}
                  {canSubmitImplementation && '提交基础成果'}
                  {canSubmitFinal && '提交最终成果'}
                </button>
              ) : (
                <div className="submit-form">
                  <h3 className="section-title">
                    {submitType === 'design' && '提交设计方案'}
                    {submitType === 'implementation' && '提交基础成果'}
                    {submitType === 'final' && '提交最终成果'}
                  </h3>
                  <textarea
                    className="form-textarea"
                    placeholder="请详细描述你的工作内容、成果和相关链接..."
                    value={submitContent}
                    onChange={(e) => setSubmitContent(e.target.value)}
                    disabled={actionLoading}
                    rows={6}
                  />
                  <div className="submit-form-actions">
                    <button
                      className="btn btn-secondary"
                      onClick={() => {
                        setShowSubmitForm(false);
                        setSubmitContent('');
                        setSubmitType('');
                      }}
                      disabled={actionLoading}
                    >
                      取消
                    </button>
                    <button
                      className="btn btn-primary"
                      onClick={handleSubmitWork}
                      disabled={actionLoading}
                    >
                      {actionLoading ? '提交中...' : '提交'}
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Cancel Task Button (for creator) */}
          {canCancelTask && (
            <div className="detail-section">
              <button
                className="btn btn-danger"
                onClick={handleCancelTask}
                disabled={actionLoading}
              >
                {actionLoading ? '取消中...' : '取消任务'}
              </button>
            </div>
          )}

          {/* Quit Task Button (for executor) */}
          {canQuitTask && (
            <div className="detail-section">
              <button
                className="btn btn-danger"
                onClick={handleQuitTask}
                disabled={actionLoading}
              >
                {actionLoading ? '退出中...' : '中途退出'}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default TaskDetailModal;
