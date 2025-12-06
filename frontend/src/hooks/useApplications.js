import { useState, useEffect } from 'react';
import { getApplications } from '../services/applicationsService';

const useApplications = (namespace = 'all') => {
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchApplications = async () => {
      try {
        setLoading(true);
        const data = await getApplications(namespace);
        setApplications(data || []);
        setError(null);
      } catch (err) {
        console.error('Error fetching applications:', err);
        setError('Failed to fetch applications');
      } finally {
        setLoading(false);
      }
    };

    fetchApplications();

    // Refresh applications every 30 seconds
    const intervalId = setInterval(fetchApplications, 30000);

    return () => clearInterval(intervalId);
  }, [namespace]);

  return { applications, loading, error };
};

export default useApplications;
