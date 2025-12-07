/**
 * Environment configuration helper
 * - In development: reads from process.env (build-time)
 * - In production (Docker): reads from window._env_ (runtime)
 */
const getEnv = (key, defaultValue = '') => {
  // Runtime config (Docker container)
  if (window._env_ && window._env_[key] !== undefined && window._env_[key] !== '') {
    return window._env_[key];
  }

  // Build-time config (development)
  if (process.env[key] !== undefined && process.env[key] !== '') {
    return process.env[key];
  }

  return defaultValue;
};

export const config = {
  API_URL: '/api',
  WS_URL: '/ws',
  REFRESH_INTERVAL: parseInt(getEnv('REACT_APP_REFRESH_INTERVAL', '5000'), 10),
  VERSION: getEnv('REACT_APP_VERSION', '1.0.0'),
  LOGO_URL: getEnv('REACT_APP_LOGO_URL', ''),
  GITHUB_URL: getEnv('REACT_APP_GITHUB_URL', ''),
  LINKEDIN_URL: getEnv('REACT_APP_LINKEDIN_URL', ''),
  INFRA_OS: getEnv('REACT_APP_INFRA_OS', 'Talos'),
  INFRA_STORAGE: getEnv('REACT_APP_INFRA_STORAGE', 'Longhorn/OpenEBS'),
  INFRA_CNI: getEnv('REACT_APP_INFRA_CNI', 'Cilium'),
  INFRA_LOAD_BALANCER: getEnv('REACT_APP_INFRA_LOAD_BALANCER', 'Traefik'),
  INFRA_GITOPS: getEnv('REACT_APP_INFRA_GITOPS', 'Flux CD'),
};

export default config;
