import config from '../config/env';

// API configuration
export const API_BASE_URL = config.API_URL;
export const WS_BASE_URL = config.WS_URL;
export const REFRESH_INTERVAL = config.REFRESH_INTERVAL;

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
