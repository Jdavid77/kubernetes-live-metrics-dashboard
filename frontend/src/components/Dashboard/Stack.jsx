import React from 'react';
import {
  FiServer,
  FiHardDrive,
  FiWifi,
  FiShare2,
  FiGitBranch,
  FiGithub,
  FiLinkedin,
} from 'react-icons/fi';
import config from '../../config/env';
import './Stack.css';

const Stack = () => {
  const stackItems = [
    {
      category: 'OS',
      value: config.INFRA_OS,
      icon: FiServer,
      bgColor: 'linear-gradient(135deg, #ff6b35 0%, #ff8952 100%)',
    },
    {
      category: 'Storage',
      value: config.INFRA_STORAGE,
      icon: FiHardDrive,
      bgColor: 'linear-gradient(135deg, #7c3aed 0%, #9333ea 100%)',
    },
    {
      category: 'CNI',
      value: config.INFRA_CNI,
      icon: FiWifi,
      bgColor: 'linear-gradient(135deg, #f59e0b 0%, #f97316 100%)',
    },
    {
      category: 'Load Balancing',
      value: config.INFRA_LOAD_BALANCER,
      icon: FiShare2,
      bgColor: 'linear-gradient(135deg, #06b6d4 0%, #0891b2 100%)',
    },
    {
      category: 'GitOps',
      value: config.INFRA_GITOPS,
      icon: FiGitBranch,
      bgColor: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
    },
  ];

  const githubUrl = config.GITHUB_URL;
  const linkedinUrl = config.LINKEDIN_URL;

  return (
    <div className="stack-section">
      <h1 className="main-title">Home Cluster Live Status</h1>
      <p className="main-description">
        This cluster is built on Talos OS, with initial configurations generated using Talhelper.
        Following deployment, Flux continuously monitors a git repository for changes, while
        Renovate handles automated dependency updates.
      </p>
      <div className="social-links">
        {githubUrl && (
          <a href={githubUrl} target="_blank" rel="noopener noreferrer" className="social-icon">
            <FiGithub size={24} />
          </a>
        )}
        {linkedinUrl && (
          <a href={linkedinUrl} target="_blank" rel="noopener noreferrer" className="social-icon">
            <FiLinkedin size={24} />
          </a>
        )}
      </div>
      <h2 className="section-title">Infrastructure Stack</h2>
      <div className="stack-grid">
        {stackItems.map((item, index) => {
          const Icon = item.icon;
          return (
            <div key={index} className="stack-card" style={{ background: item.bgColor }}>
              <div className="stack-icon">
                <Icon size={32} />
              </div>
              <div className="stack-content">
                <div className="stack-category">{item.category}</div>
                <div className="stack-value">{item.value}</div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default Stack;
