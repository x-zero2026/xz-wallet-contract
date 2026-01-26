// Authentication utilities
export const getToken = () => {
  return localStorage.getItem('token');
};

export const setToken = (token) => {
  localStorage.setItem('token', token);
};

export const removeToken = () => {
  localStorage.removeItem('token');
};

export const getUserInfo = () => {
  const userStr = localStorage.getItem('userInfo');
  return userStr ? JSON.parse(userStr) : null;
};

export const setUserInfo = (userInfo) => {
  localStorage.setItem('userInfo', JSON.stringify(userInfo));
};

export const removeUserInfo = () => {
  localStorage.removeItem('userInfo');
};

export const isAuthenticated = () => {
  return !!getToken();
};

export const logout = () => {
  removeToken();
  removeUserInfo();
};
