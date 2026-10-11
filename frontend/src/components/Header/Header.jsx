import React from 'react';
import { NavLink } from 'react-router-dom';
import { FiActivity, FiGrid, FiLayers } from 'react-icons/fi';
import config from '../../config/env';
import './Header.css';

const Header = ({ isConnected }) => {
  const logoUrl = config.LOGO_URL;
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
      <nav className="header-nav">
        <div className="container nav-tabs">
          <NavLink to="/" end className={({ isActive }) => `nav-tab ${isActive ? 'active' : ''}`}>
            <FiGrid size={16} /> Dashboard
          </NavLink>
          <NavLink to="/architecture" className={({ isActive }) => `nav-tab ${isActive ? 'active' : ''}`}>
            <FiLayers size={16} /> Architecture
          </NavLink>
        </div>
      </nav>
    </header>
  );
};

export default Header;
