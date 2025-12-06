import React from 'react';
import Stack from './Stack';
import ClusterOverview from './ClusterOverview';
import NodesSection from './NodesSection';
import useMetrics from '../../hooks/useMetrics';
import './Dashboard.css';

const Dashboard = () => {
  const { metrics } = useMetrics();

  return (
    <div className="dashboard">
      <div className="container">
        <Stack />
        <ClusterOverview metrics={metrics} />
        <NodesSection />
      </div>
    </div>
  );
};

export default Dashboard;
