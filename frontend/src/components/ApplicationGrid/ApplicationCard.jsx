import React from 'react';
import { FiCheckCircle, FiAlertCircle, FiXCircle } from 'react-icons/fi';
import { STATUS_COLORS } from '../../utils/constants';
import './ApplicationCard.css';

const ApplicationCard = ({ application }) => {
  const getStatusIcon = () => {
    switch (application.status) {
      case 'Healthy':
        return <FiCheckCircle />;
      case 'Degraded':
        return <FiAlertCircle />;
      case 'Down':
        return <FiXCircle />;
      default:
        return <FiAlertCircle />;
    }
  };

  return (
    <div className="application-card">
      <div className="app-header">
        <h3 className="app-name">{application.name}</h3>
        <div className="app-status" style={{ color: STATUS_COLORS[application.status] }}>
          {getStatusIcon()}
          <span>{application.status}</span>
        </div>
      </div>

      <div className="app-details">
        <div className="app-detail-item">
          <span className="detail-label">Namespace:</span>
          <span className="detail-value">{application.namespace}</span>
        </div>
        <div className="app-detail-item">
          <span className="detail-label">Pods:</span>
          <span className="detail-value">{application.pod_count}</span>
        </div>
      </div>

      {application.technology && application.technology.length > 0 && (
        <div className="app-tech">
          {application.technology.map((tech, index) => (
            <span key={index} className="tech-badge">
              {tech}
            </span>
          ))}
        </div>
      )}
    </div>
  );
};

export default ApplicationCard;
