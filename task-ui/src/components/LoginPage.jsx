import { useState } from 'react';
import { setToken, setUserInfo } from '../utils/auth';
import './LoginPage.css';

function LoginPage({ onLogin }) {
  const [token, setTokenInput] = useState('');
  const [userInfoInput, setUserInfoInput] = useState('');
  const [error, setError] = useState('');

  const handleLogin = () => {
    try {
      if (!token.trim()) {
        setError('请输入 JWT Token');
        return;
      }

      // Parse user info
      let userInfo;
      if (userInfoInput.trim()) {
        userInfo = JSON.parse(userInfoInput);
      } else {
        // Default admin user info for testing
        userInfo = {
          did: '0x3070deb1c17432b094d30509ccbfd598fb2793435efdca9273dfbc558bc040ca',
          username: 'admin',
          eth_address: '0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA',
          email: 'x-zero2026@outlook.com',
          project_id: '00000000-0000-0000-0000-000000000000',
        };
      }

      setToken(token);
      setUserInfo(userInfo);
      onLogin();
    } catch (err) {
      setError('用户信息格式错误: ' + err.message);
    }
  };

  const handleQuickLogin = () => {
    // For testing: use admin token
    const adminToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkaWQiOiIweDMwNzBkZWIxYzE3NDMyYjA5NGQzMDUwOWNjYmZkNTk4ZmIyNzkzNDM1ZWZkY2E5MjczZGZiYzU1OGJjMDQwY2EiLCJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzY5NTEzNTg2LCJpYXQiOjE3Njk0MjcxODZ9.xFST95uXdhMYbCy7p_B9BVR1FRqPzeDEDHPxNY8Rv4Q';
    const adminInfo = {
      did: '0x3070deb1c17432b094d30509ccbfd598fb2793435efdca9273dfbc558bc040ca',
      username: 'admin',
      eth_address: '0x8766E5c6c7311519187328Ebe31fD63C4b88A9cA',
      email: 'x-zero2026@outlook.com',
      project_id: '00000000-0000-0000-0000-000000000000',
    };

    setToken(adminToken);
    setUserInfo(adminInfo);
    onLogin();
  };

  return (
    <div className="login-page">
      <div className="login-card">
        <h1 className="login-title">XZ 任务中心</h1>
        <p className="login-subtitle">请登录以继续</p>

        {error && <div className="error">{error}</div>}

        <div className="form-group">
          <label className="form-label">JWT Token</label>
          <textarea
            className="form-textarea"
            placeholder="粘贴你的 JWT Token"
            value={token}
            onChange={(e) => setTokenInput(e.target.value)}
            rows={4}
          />
        </div>

        <div className="form-group">
          <label className="form-label">用户信息 (JSON, 可选)</label>
          <textarea
            className="form-textarea"
            placeholder='{"did": "0x...", "username": "admin", "eth_address": "0x...", "project_id": "uuid"}'
            value={userInfoInput}
            onChange={(e) => setUserInfoInput(e.target.value)}
            rows={4}
          />
          <small style={{ color: '#999', fontSize: '0.85rem' }}>
            留空将使用默认 admin 用户信息
          </small>
        </div>

        <button className="btn btn-primary btn-block" onClick={handleLogin}>
          登录
        </button>

        <div style={{ marginTop: '1rem', textAlign: 'center' }}>
          <button className="btn btn-secondary" onClick={handleQuickLogin}>
            快速登录（测试用）
          </button>
        </div>

        <div className="login-help">
          <p>提示：</p>
          <ul>
            <li>从 DID Login 系统获取 JWT Token</li>
            <li>或使用快速登录进行测试</li>
          </ul>
        </div>
      </div>
    </div>
  );
}

export default LoginPage;
