import { formatDistanceToNow } from 'date-fns';

// Format timestamp to relative time
export const formatRelativeTime = (timestamp) => {
  if (!timestamp) return 'N/A';
  try {
    return formatDistanceToNow(new Date(timestamp), { addSuffix: true });
  } catch (error) {
    return 'Invalid date';
  }
};

// Format number with commas
export const formatNumber = (num) => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

// Format percentage
export const formatPercentage = (value) => {
  if (typeof value === 'string' && value.endsWith('%')) {
    return value;
  }
  return `${value}%`;
};

// Format memory
export const formatMemory = (value) => {
  if (typeof value === 'string' && (value.endsWith('GB') || value.endsWith('MB'))) {
    return value;
  }
  return `${value}GB`;
};
