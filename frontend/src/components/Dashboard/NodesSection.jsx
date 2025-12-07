import React, { useState, useEffect } from 'react';
import { FiChevronDown, FiChevronUp, FiCheckCircle, FiXCircle, FiServer } from 'react-icons/fi';
import { getNodes } from '../../services/nodesService';
import LoadingSpinner from '../Common/LoadingSpinner';
import './NodesSection.css';

const NodesSection = () => {
  const [nodes, setNodes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isExpanded, setIsExpanded] = useState(() => {
    const saved = localStorage.getItem('nodes-section-expanded');
    return saved ? JSON.parse(saved) : false;
  });

  useEffect(() => {
    const fetchNodes = async () => {
      try {
        setLoading(true);
        const data = await getNodes();
        setNodes(data || []);
        setError(null);
      } catch (err) {
        console.error('Error fetching nodes:', err);
        setError('Failed to fetch nodes');
      } finally {
        setLoading(false);
      }
    };

    fetchNodes();
    const interval = setInterval(fetchNodes, 30000); // Refresh every 30 seconds
    return () => clearInterval(interval);
  }, []);

  const toggleExpanded = () => {
    const newExpanded = !isExpanded;
    setIsExpanded(newExpanded);
    localStorage.setItem('nodes-section-expanded', JSON.stringify(newExpanded));
  };

  const getUsageColor = (percentage) => {
    if (percentage < 60) return '#10b981'; // Green
    if (percentage < 80) return '#f59e0b'; // Yellow
    return '#ef4444'; // Red
  };

  return (
    <div className="nodes-section">
      <div className="nodes-header" onClick={toggleExpanded}>
        <h2 className="nodes-header-title">Cluster Nodes</h2>
        <div className="nodes-header-right">
          {!loading && <span className="nodes-count">{nodes.length} node{nodes.length !== 1 ? 's' : ''}</span>}
          {isExpanded ? <FiChevronUp className="chevron-icon" size={24} /> : <FiChevronDown className="chevron-icon" size={24} />}
        </div>
      </div>

      {isExpanded && (
        <div className="nodes-content">
          {loading ? (
            <LoadingSpinner />
          ) : error ? (
            <div className="text-center text-muted">{error}</div>
          ) : nodes.length === 0 ? (
            <div className="text-center text-muted">No nodes found</div>
          ) : (
            <div className="nodes-grid">
              {nodes.map((node, index) => {
                const isReady = node.status === 'Ready';
                const cpuColor = getUsageColor(node.cpu_usage);
                const memColor = getUsageColor(node.memory_usage);

                return (
                  <div key={index} className="node-card">
                    <div className="node-header-row">
                      <div className="node-name-container">
                        <div className="node-icon">
                          <FiServer size={20} />
                        </div>
                        <div className="node-name">
                          {node.name}
                          <span className={`node-role-badge ${node.role === 'control-plane' ? 'control-plane' : 'worker'}`}>
                            {node.role === 'control-plane' ? 'Control Plane' : 'Worker'}
                          </span>
                        </div>
                      </div>
                      <div className={`node-status-badge ${isReady ? 'ready' : 'not-ready'}`}>
                        {isReady ? (
                          <>
                            <FiCheckCircle size={16} />
                            <span>Ready</span>
                          </>
                        ) : (
                          <>
                            <FiXCircle size={16} />
                            <span>NotReady</span>
                          </>
                        )}
                      </div>
                    </div>

                    <div className="node-metrics">
                      <div className="metric-row">
                        <div className="metric-label">CPU Usage</div>
                        <div className="metric-value">{node.cpu_usage.toFixed(1)}%</div>
                      </div>
                      <div className="progress-bar-container">
                        <div
                          className="progress-bar-fill"
                          style={{
                            width: `${Math.min(node.cpu_usage, 100)}%`,
                            backgroundColor: cpuColor
                          }}
                        />
                      </div>

                      <div className="metric-row">
                        <div className="metric-label">Memory Usage</div>
                        <div className="metric-value">{node.memory_usage.toFixed(1)}%</div>
                      </div>
                      <div className="progress-bar-container">
                        <div
                          className="progress-bar-fill"
                          style={{
                            width: `${Math.min(node.memory_usage, 100)}%`,
                            backgroundColor: memColor
                          }}
                        />
                      </div>

                      <div className="node-footer">
                        <div className="node-info-item">
                          <span className="info-label">CPU:</span>
                          <span className="info-value">{node.cpu_capacity}</span>
                        </div>
                        <div className="node-info-item">
                          <span className="info-label">Memory:</span>
                          <span className="info-value">{node.memory_capacity}</span>
                        </div>
                        <div className="node-info-item">
                          <span className="info-label">Pods:</span>
                          <span className="info-value">{node.pod_count}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default NodesSection;
