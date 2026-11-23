import yaml from 'js-yaml';

/**
 * Runtime configuration, loaded from /config/config.yaml before the app renders
 * (see index.js). Override per deployment by mounting a ConfigMap over
 * /usr/share/nginx/html/config/. Values missing from the file keep these defaults.
 */
const CONFIG_URL = `${process.env.PUBLIC_URL}/config/config.yaml`;

export const config = {
  API_URL: '/api',
  WS_URL: '/ws',
  REFRESH_INTERVAL: 5000,
  LOGO_URL: '',
  GITHUB_URL: '',
  LINKEDIN_URL: '',
  INFRA_OS: 'Talos',
  INFRA_STORAGE: 'Longhorn/OpenEBS',
  INFRA_CNI: 'Cilium',
  INFRA_LOAD_BALANCER: 'Traefik',
  INFRA_GITOPS: 'Flux CD',
  ARCHITECTURE_TOPICS: null,
};

// Never rejects: on failure the app renders with the defaults above.
export const loadConfig = async () => {
  try {
    const res = await fetch(CONFIG_URL, { cache: 'no-cache' });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = yaml.load(await res.text());
    // nginx falls back to index.html for missing files, so check the shape.
    if (typeof data !== 'object' || data === null) throw new Error('not a YAML mapping');

    const { links = {}, stack = {}, architecture = {} } = data;
    const refreshInterval = Number(data.refreshInterval);
    const values = {
      REFRESH_INTERVAL: refreshInterval > 0 ? refreshInterval : undefined,
      LOGO_URL: data.logoUrl,
      GITHUB_URL: links.github,
      LINKEDIN_URL: links.linkedin,
      INFRA_OS: stack.os,
      INFRA_STORAGE: stack.storage,
      INFRA_CNI: stack.cni,
      INFRA_LOAD_BALANCER: stack.loadBalancer,
      INFRA_GITOPS: stack.gitops,
      ARCHITECTURE_TOPICS: Array.isArray(architecture.topics) ? architecture.topics : undefined,
    };
    Object.entries(values).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') config[key] = value;
    });
  } catch (err) {
    console.error('Error loading config, using defaults:', err);
  }
};

export default config;
