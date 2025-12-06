import React from 'react';
import { FiCheckCircle, FiAlertCircle, FiXCircle } from 'react-icons/fi';
import { formatRelativeTime } from '../../utils/formatters';
import LoadingSpinner from '../Common/LoadingSpinner';
import './ResourceTable.css';

const DeploymentsTable = ({ applications, loading }) => {
  if (loading) {
    return <LoadingSpinner />;
  }

  const getStatusIcon = (status) => {
    switch (status) {
      case 'Healthy':
        return <FiCheckCircle className="icon-success" size={18} />;
      case 'Degraded':
        return <FiAlertCircle className="icon-warning" size={18} />;
      case 'Down':
        return <FiXCircle className="icon-danger" size={18} />;
      default:
        return null;
    }
  };

  return (
    <div className="resource-table-container">
      <div className="table-header">
        <h2 className="section-title">Deployments</h2>
        <div className="table-stats">
          Total: {applications?.length || 0} deployments
        </div>
      </div>

      <div className="table-wrapper">
        <table className="resource-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Namespace</th>
              <th>Status</th>
              <th>Pods</th>
              <th>Age</th>
              <th>Image</th>
            </tr>
          </thead>
          <tbody>
            {applications?.map((app, index) => (
              <tr key={index}>
                <td className="deployment-name">{app.name}</td>
                <td>
                  <span className="namespace-badge">{app.namespace}</span>
                </td>
                <td>
                  <div className="status-cell">
                    {getStatusIcon(app.status)}
                    <span className={`status-badge ${app.status.toLowerCase()}`}>
                      {app.status}
                    </span>
                  </div>
                </td>
                <td>
                  <span className="pod-count">{app.pod_count}</span>
                </td>
                <td className="age-cell">
                  {formatRelativeTime(app.created_at)}
                </td>
                <td className="image-cell">
                  <code>{app.image}</code>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default DeploymentsTable;
