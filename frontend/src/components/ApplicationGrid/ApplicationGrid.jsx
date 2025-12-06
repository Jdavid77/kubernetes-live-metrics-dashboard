import React from 'react';
import ApplicationCard from './ApplicationCard';
import LoadingSpinner from '../Common/LoadingSpinner';
import './ApplicationGrid.css';

const ApplicationGrid = ({ applications, loading }) => {
  if (loading) {
    return <LoadingSpinner />;
  }

  if (!applications || applications.length === 0) {
    return (
      <div className="no-applications">
        <p>No applications found</p>
      </div>
    );
  }

  return (
    <div className="applications-section">
      <h2 className="section-title">Applications</h2>
      <div className="application-grid">
        {applications.map((app, index) => (
          <ApplicationCard key={index} application={app} />
        ))}
      </div>
    </div>
  );
};

export default ApplicationGrid;
