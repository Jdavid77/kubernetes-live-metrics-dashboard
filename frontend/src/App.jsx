import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Header from './components/Header/Header';
import Dashboard from './components/Dashboard/Dashboard';
import Architecture from './components/Architecture/Architecture';
import ErrorBoundary from './components/Common/ErrorBoundary';
import useMetrics from './hooks/useMetrics';
import './App.css';

function App() {
  const { isConnected } = useMetrics();

  return (
    <ErrorBoundary>
      <BrowserRouter>
        <div className="App">
          <Header isConnected={isConnected} />
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/architecture" element={<Architecture />} />
          </Routes>
        </div>
      </BrowserRouter>
    </ErrorBoundary>
  );
}

export default App;
