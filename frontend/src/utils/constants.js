// API configuration
export const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
export const WS_BASE_URL = process.env.REACT_APP_WS_URL || 'ws://localhost:8080';
export const REFRESH_INTERVAL = parseInt(process.env.REACT_APP_REFRESH_INTERVAL || '5000', 10);

// Status colors
export const STATUS_COLORS = {
  Healthy: '#10b981',
  Degraded: '#f59e0b',
  Down: '#ef4444',
};

// Application status
export const APP_STATUS = {
  HEALTHY: 'Healthy',
  DEGRADED: 'Degraded',
  DOWN: 'Down',
};
