import { useState, useEffect } from 'react';
import useWebSocket from './useWebSocket';
import { getMetrics } from '../services/metricsService';
import { REFRESH_INTERVAL } from '../utils/constants';

const useMetrics = () => {
  const [metrics, setMetrics] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { data: wsData, error: wsError, isConnected } = useWebSocket('/ws/metrics');

  // Use WebSocket data when available
  useEffect(() => {
    if (wsData) {
      setMetrics(wsData);
      setLoading(false);
      setError(null);
    }
  }, [wsData]);

  // Fallback to polling if WebSocket fails
  useEffect(() => {
    let intervalId;

    if (wsError || !isConnected) {
      // Fetch initial data
      const fetchMetrics = async () => {
        try {
          const data = await getMetrics();
          setMetrics(data);
          setLoading(false);
          setError(null);
        } catch (err) {
          console.error('Error fetching metrics:', err);
          setError('Failed to fetch metrics');
          setLoading(false);
        }
      };

      fetchMetrics();

      // Poll every REFRESH_INTERVAL
      intervalId = setInterval(fetchMetrics, REFRESH_INTERVAL);
    }

    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [wsError, isConnected]);

  return { metrics, loading, error, isConnected };
};

export default useMetrics;
