import React from 'react';
import { FiActivity } from 'react-icons/fi';
import config from '../../config/env';
import './Header.css';

const Header = ({ isConnected }) => {
  const logoUrl = config.LOGO_URL;
  const version = config.VERSION;

  return (
    <header className="header">
      <div className="header-top">
        <div className="container">
          <div className="header-brand">
            <div className={`k8s-logo ${logoUrl ? 'has-image' : ''}`}>
              {logoUrl ? (
                <img src={logoUrl} alt="Logo" className="logo-image" />
              ) : (
                <FiActivity size={32} />
              )}
            </div>
            <div>
              <h1>Kubernetes Dashboard</h1>
              <span className="version">v{version}</span>
            </div>
          </div>

          <div className="header-controls">
            <div className="status-indicator-container">
              <span className="status-text">API Status</span>
              <span className={`status-dot ${isConnected ? 'connected' : 'disconnected'}`}></span>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
