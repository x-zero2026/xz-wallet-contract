import { useState, useEffect } from 'react';
import { getUserInfo } from './utils/auth';
import { getBalance, listProjects } from './api';
import TaskList from './components/TaskList';
import CreateTaskModal from './components/CreateTaskModal';
import TaskDetailModal from './components/TaskDetailModal';
import StarfieldBackground from './components/StarfieldBackground';
import './App.css';

function App() {
  const [activeTab, setActiveTab] = useState('my-published');
  const [userInfo, setUserInfo] = useState(null);
  const [balance, setBalance] = useState('0');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [selectedTaskId, setSelectedTaskId] = useState(null);
  const [refreshKey, setRefreshKey] = useState(0);
  const [projects, setProjects] = useState([]);
  const [selectedProject, setSelectedProject] = useState(null);
  const [loading, setLoading] = useState(false); // 改为 false，不显示初始加载

  useEffect(() => {
    const initializeApp = async () => {
      // Check if token is in URL (from DID login redirect)
      const urlParams = new URLSearchParams(window.location.search);
      const tokenFromUrl = urlParams.get('token');
      
      if (tokenFromUrl) {
        console.log('✅ Token received from URL');
        // Save token to localStorage
        localStorage.setItem('token', tokenFromUrl);
        
        // Parse JWT to get user info
        try {
          const payload = JSON.parse(atob(tokenFromUrl.split('.')[1]));
          console.log('Token payload:', payload);
          
          // Save basic user info from token
          // Note: eth_address needs to be queried from database or will be fetched when needed
          const userInfo = {
            did: payload.did,
            username: payload.username,
            eth_address: '', // Will be populated when fetching balance
          };
          localStorage.setItem('userInfo', JSON.stringify(userInfo));
          console.log('✅ User info saved:', userInfo);
          
          // Remove token from URL for security
          window.history.replaceState({}, document.title, window.location.pathname);
          
          // Reload to apply changes
          window.location.reload();
          return;
        } catch (err) {
          console.error('Failed to parse token:', err);
        }
      }
      
      const user = getUserInfo();
      if (!user) {
        // 未登录，不需要设置 loading
        return;
      }
      setUserInfo(user);
      
      // Load projects
      await loadProjects();
      
      // Fetch eth_address if not present
      if (!user.eth_address && user.did) {
        await fetchUserEthAddress(user.did);
      } else if (user.eth_address) {
        await loadBalance(user.eth_address);
      }
    };
    
    initializeApp();
  }, []);

  const fetchUserEthAddress = async (did) => {
    try {
      // Fetch from DID Login API
      const DID_LOGIN_API_URL = import.meta.env.VITE_DID_LOGIN_API_URL;
      const token = localStorage.getItem('token');
      
      const response = await fetch(`${DID_LOGIN_API_URL}/api/user/profile`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });
      
      if (response.ok) {
        const data = await response.json();
        const ethAddress = data.data?.eth_address;
        
        if (ethAddress) {
          const updatedUserInfo = { ...getUserInfo(), eth_address: ethAddress };
          localStorage.setItem('userInfo', JSON.stringify(updatedUserInfo));
          setUserInfo(updatedUserInfo);
          loadBalance(ethAddress);
        }
      }
    } catch (err) {
      console.error('Failed to fetch eth_address:', err);
    }
  };

  const loadBalance = async (address) => {
    try {
      const response = await getBalance(address);
      const rawBalance = response.data.data.xzt_balance;
      // Format to 2 decimal places
      setBalance(parseFloat(rawBalance).toFixed(2));
    } catch (err) {
      console.error('Load balance error:', err);
    }
  };

  const loadProjects = async () => {
    try {
      const response = await listProjects();
      const projectList = response.data.data || [];
      setProjects(projectList);
      
      // 默认选择第一个项目
      if (projectList.length > 0) {
        setSelectedProject(projectList[0]);
      }
    } catch (err) {
      console.error('Load projects error:', err);
    }
  };

  const handleRefresh = () => {
    setRefreshKey((prev) => prev + 1);
    if (userInfo) {
      loadBalance(userInfo.eth_address);
    }
  };

  const getTaskFilter = () => {
    switch (activeTab) {
      case 'available':
        // 可申请任务：状态为 bidding，且不是自己发布的，且没有投标过的
        return { 
          status: 'bidding',
          exclude_creator: true,  // 排除自己发布的任务
          exclude_bidded: true    // 排除已投标的任务
        };
      case 'my-tasks':
        // 我的任务：我已投标的任务（包括还在招标中的和已被选中的）
        return { bidder_did: userInfo?.did };
      case 'my-published':
        return { creator_did: userInfo?.did };
      default:
        return {};
    }
  };

  if (!userInfo) {
    return (
      <div className="app">
        <div className="main-content">
          <div className="card" style={{ textAlign: 'center', maxWidth: '500px', margin: '2rem auto' }}>
            <h2>欢迎使用 XZ 任务中心</h2>
            <p style={{ color: '#666', margin: '1rem 0' }}>
              请先登录以访问任务中心
            </p>
            <p style={{ color: '#999', fontSize: '0.9rem' }}>
              您需要通过 DID 登录系统进行身份验证
            </p>
            <a 
              href="https://main.d2fozf421c6ftf.amplifyapp.com" 
              className="btn btn-primary"
              style={{ marginTop: '1rem', display: 'inline-block' }}
            >
              前往登录
            </a>
          </div>
        </div>
      </div>
    );
  }

  return (
    <>
      <StarfieldBackground />
      <div className="app">
      {/* Header */}
      <header className="header">
        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
          <h1 className="header-title">XZ 任务中心</h1>
          {/* 项目选择器 */}
          {projects.length > 0 && (
            <select
              id="project-select"
              value={selectedProject?.project_id || ''}
              onChange={(e) => {
                const project = projects.find(p => p.project_id === e.target.value);
                setSelectedProject(project);
              }}
              style={{
                padding: '0.5rem',
                borderRadius: '4px',
                border: '1px solid #ddd',
                cursor: 'pointer',
              }}
            >
              {projects.map((project) => (
                <option key={project.project_id} value={project.project_id}>
                  {project.project_name}
                </option>
              ))}
            </select>
          )}
        </div>
        <div className="header-user">
          <div className="user-info">
            <div className="user-name">{userInfo.username}</div>
            <div className="user-balance">余额: {balance} XZT</div>
            {userInfo.eth_address && (
              <div className="user-address">
                钱包: {userInfo.eth_address}
              </div>
            )}
          </div>
        </div>
      </header>

      {/* Navigation */}
      <nav className="nav">
        <ul className="nav-tabs">
          <li
            className={`nav-tab ${activeTab === 'my-published' ? 'active' : ''}`}
            onClick={() => setActiveTab('my-published')}
          >
            我发布的
          </li>
          <li
            className={`nav-tab ${activeTab === 'available' ? 'active' : ''}`}
            onClick={() => setActiveTab('available')}
          >
            可申请任务
          </li>
          <li
            className={`nav-tab ${activeTab === 'my-tasks' ? 'active' : ''}`}
            onClick={() => setActiveTab('my-tasks')}
          >
            我的任务
          </li>
        </ul>
      </nav>

      {/* Main Content */}
      <main className="main-content">
        {activeTab === 'my-published' && (
          <div style={{ marginBottom: '1.5rem' }}>
            <button
              className="btn btn-primary"
              onClick={() => setShowCreateModal(true)}
            >
              + 发布新任务
            </button>
          </div>
        )}

        <TaskList
          key={`${activeTab}-${refreshKey}`}
          filter={getTaskFilter()}
          onTaskClick={(task) => setSelectedTaskId(task.task_id)}
        />
      </main>

      {/* Modals */}
      {showCreateModal && (
        <CreateTaskModal
          onClose={() => setShowCreateModal(false)}
          onSuccess={handleRefresh}
          selectedProject={selectedProject}
        />
      )}

      {selectedTaskId && (
        <TaskDetailModal
          taskId={selectedTaskId}
          onClose={() => setSelectedTaskId(null)}
          onUpdate={handleRefresh}
          onSwitchTab={(tab) => setActiveTab(tab)}
        />
      )}
    </div>
    </>
  );
}

export default App;
