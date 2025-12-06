import React from 'react';
import Header from './components/Header/Header';
import Dashboard from './components/Dashboard/Dashboard';
import ErrorBoundary from './components/Common/ErrorBoundary';
import useMetrics from './hooks/useMetrics';
import './App.css';

function App() {
  const { isConnected } = useMetrics();

  return (
    <ErrorBoundary>
      <div className="App">
        <Header isConnected={isConnected} />
        <Dashboard />
      </div>
    </ErrorBoundary>
  );
}

export default App;
