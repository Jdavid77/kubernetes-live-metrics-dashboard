import React from 'react';
import { FiCpu, FiDatabase, FiBox, FiGlobe, FiServer } from 'react-icons/fi';
import './ClusterOverview.css';

const ClusterOverview = ({ metrics }) => {
  const healthMetrics = [
    {
      title: 'CPU Usage',
      value: metrics?.cpu_usage || '0 cores',
      icon: FiCpu,
      color: '#3b82f6',
    },
    {
      title: 'Memory Usage',
      value: metrics?.memory_usage || '0GB',
      icon: FiDatabase,
      color: '#8b5cf6',
    },
  ];

  const resources = [
    {
      label: 'Running Pods',
      value: metrics?.pods_running || 0,
      icon: FiBox,
      color: '#10b981',
    },
    {
      label: 'Services',
      value: metrics?.services_count || 0,
      icon: FiGlobe,
      color: '#f59e0b',
    },
    {
      label: 'HTTP Routes',
      value: metrics?.httproute_count || 0,
      icon: FiGlobe,
      color: '#8b5cf6',
    },
    {
      label: 'Gateways',
      value: metrics?.gateway_count || 0,
      icon: FiServer,
      color: '#ec4899',
    },
  ];

  return (
    <div className="cluster-overview">
      <h2 className="section-title">Cluster Overview</h2>
      <div className="overview-grid">
        <div className="health-section">
          {healthMetrics.map((metric, index) => {
            const Icon = metric.icon;
            return (
              <div key={index} className="health-card" style={{ '--card-color': metric.color }}>
                <div className="health-icon" style={{ backgroundColor: `${metric.color}20`, color: metric.color }}>
                  <Icon size={28} />
                </div>
                <div className="health-content">
                  <div className="health-title">{metric.title}</div>
                  <div className="health-value">{metric.value}</div>
                </div>
              </div>
            );
          })}
        </div>

        <div className="resources-section">
          {resources.map((resource, index) => {
            const Icon = resource.icon;
            return (
              <div key={index} className="resource-item">
                <div className="resource-icon-small" style={{ backgroundColor: `${resource.color}20`, color: resource.color }}>
                  <Icon size={20} />
                </div>
                <div className="resource-details">
                  <div className="resource-count">{resource.value}</div>
                  <div className="resource-name">{resource.label}</div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

export default ClusterOverview;
