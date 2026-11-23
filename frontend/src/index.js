import React from 'react';
import ReactDOM from 'react-dom/client';
import './styles/variables.css';
import './styles/globals.css';
import App from './App';
import { loadConfig } from './config/env';

const root = ReactDOM.createRoot(document.getElementById('root'));
loadConfig().then(() => {
  root.render(
    <React.StrictMode>
      <App />
    </React.StrictMode>
  );
});
