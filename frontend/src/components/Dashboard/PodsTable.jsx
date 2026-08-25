import React, { useState, useEffect } from 'react';
import { FiCheckCircle, FiXCircle } from 'react-icons/fi';
import { getPods } from '../../services/podsService';
import LoadingSpinner from '../Common/LoadingSpinner';
import './ResourceTable.css';

const PodsTable = ({ namespace }) => {
  const [pods, setPods] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchPods = async () => {
      try {
        setLoading(true);
        const data = await getPods(namespace);
        setPods(data || []);
        setError(null);
      } catch (err) {
        console.error('Error fetching pods:', err);
        setError('Failed to fetch pods');
      } finally {
        setLoading(false);
      }
    };

    fetchPods();
  }, [namespace]);

  if (loading) {
    return <LoadingSpinner />;
  }

  if (error) {
    return <div className="text-center text-muted">{error}</div>;
  }

  return (
    <div className="resource-table-container">
      <div className="table-header">
        <h2 className="section-title">Pods</h2>
        <div className="table-stats">Total: {pods.length} pods</div>
      </div>

      <div className="table-wrapper">
        <table className="resource-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Namespace</th>
              <th>Status</th>
              <th>Ready</th>
              <th>Restarts</th>
              <th>Node</th>
              <th>IP</th>
              <th>Age</th>
            </tr>
          </thead>
          <tbody>
            {pods.map((pod, index) => {
              // Parse ready status (e.g., "1/1" or "0/1")
              const [ready, total] = pod.ready.split('/').map(Number);
              const isReady = ready === total && total > 0;

              return (
                <tr key={index}>
                  <td className="pod-name">{pod.name}</td>
                  <td>
                    <span className="namespace-badge">{pod.namespace}</span>
                  </td>
                  <td>
                    <span className={`status-badge ${pod.status.toLowerCase()}`}>{pod.status}</span>
                  </td>
                  <td>
                    <div className="status-cell">
                      {isReady ? (
                        <FiCheckCircle className="icon-success" size={18} />
                      ) : (
                        <FiXCircle className="icon-danger" size={18} />
                      )}
                      <span>{pod.ready}</span>
                    </div>
                  </td>
                  <td className="age-cell">{pod.restarts}</td>
                  <td className="age-cell">{pod.node || '-'}</td>
                  <td className="ip-cell">{pod.ip || '-'}</td>
                  <td className="age-cell">{pod.age}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default PodsTable;
