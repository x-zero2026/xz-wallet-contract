import axios from 'axios';
import { getToken } from '../utils/auth';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// ============================================
// Wallet APIs
// ============================================
export const getBalance = (address) => {
  return api.get(`/wallet/balance?address=${address}`);
};

// ============================================
// Project APIs (from DID Login)
// ============================================
const DID_LOGIN_API_URL = import.meta.env.VITE_DID_LOGIN_API_URL;

export const listProjects = () => {
  const token = getToken();
  return axios.get(`${DID_LOGIN_API_URL}/api/projects`, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });
};

// ============================================
// Profession Tags APIs
// ============================================
const BUSINESS_CONSULTANT_API_URL = import.meta.env.VITE_BUSINESS_CONSULTANT_API_URL || 'https://r3jwp815n9.execute-api.us-east-1.amazonaws.com/Prod';

export const identifyProfessionTags = (taskDescription) => {
  const token = getToken();
  return axios.post(
    `${BUSINESS_CONSULTANT_API_URL}/identify-profession-tags`,
    { task_description: taskDescription },
    {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    }
  );
};

export const recommendUsers = (professionTags) => {
  const token = getToken();
  return axios.post(
    `${DID_LOGIN_API_URL}/api/recommend-users`,
    { profession_tags: professionTags },
    {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    }
  );
};

// ============================================
// Task APIs
// ============================================
export const createTask = (taskData) => {
  return api.post('/tasks', taskData);
};

export const listTasks = (params) => {
  const queryString = new URLSearchParams(params).toString();
  return api.get(`/tasks?${queryString}`);
};

export const getTask = (taskId) => {
  return api.get(`/tasks/${taskId}`);
};

export const bidTask = (taskId, bidData) => {
  return api.post(`/tasks/${taskId}/bid`, bidData);
};

export const selectBidder = (taskId, bidderData) => {
  return api.post(`/tasks/${taskId}/select-bidder`, bidderData);
};

export const approveWork = (taskId, approvalData) => {
  return api.post(`/tasks/${taskId}/approve`, approvalData);
};

export const cancelTask = (taskId) => {
  return api.post(`/tasks/${taskId}/cancel`);
};

export const submitWork = (taskId, submissionData) => {
  return api.post(`/tasks/${taskId}/submit`, submissionData);
};

// ============================================
// Task Status Constants
// ============================================
export const TASK_STATUS = {
  PENDING: 'pending',
  BIDDING: 'bidding',
  ACCEPTED: 'accepted',
  DESIGN_SUBMITTED: 'design_submitted',
  DESIGN_APPROVED: 'design_approved',
  IMPLEMENTATION_SUBMITTED: 'implementation_submitted',
  IMPLEMENTATION_APPROVED: 'implementation_approved',
  FINAL_SUBMITTED: 'final_submitted',
  COMPLETED: 'completed',
  CANCELLED: 'cancelled',
};

export const TASK_STATUS_LABELS = {
  [TASK_STATUS.PENDING]: '待发布',
  [TASK_STATUS.BIDDING]: '招标中',
  [TASK_STATUS.ACCEPTED]: '已接受',
  [TASK_STATUS.DESIGN_SUBMITTED]: '设计已提交',
  [TASK_STATUS.DESIGN_APPROVED]: '设计已批准',
  [TASK_STATUS.IMPLEMENTATION_SUBMITTED]: '实现已提交',
  [TASK_STATUS.IMPLEMENTATION_APPROVED]: '实现已批准',
  [TASK_STATUS.FINAL_SUBMITTED]: '最终成果已提交',
  [TASK_STATUS.COMPLETED]: '已完成',
  [TASK_STATUS.CANCELLED]: '已取消',
};

export const VISIBILITY = {
  PROJECT: 'project',
  GLOBAL: 'global',
};

export const VISIBILITY_LABELS = {
  [VISIBILITY.PROJECT]: '项目内可见',
  [VISIBILITY.GLOBAL]: '全局可见',
};

export default api;
