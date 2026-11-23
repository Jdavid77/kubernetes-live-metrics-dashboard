import {
  FiActivity,
  FiBox,
  FiCloud,
  FiCpu,
  FiDatabase,
  FiGitBranch,
  FiGlobe,
  FiHardDrive,
  FiKey,
  FiLayers,
  FiLock,
  FiServer,
  FiShield,
} from 'react-icons/fi';

// Names accepted by a topic's `icon` field in config.yaml.
const ICONS = {
  activity: FiActivity,
  box: FiBox,
  cloud: FiCloud,
  cpu: FiCpu,
  database: FiDatabase,
  'git-branch': FiGitBranch,
  globe: FiGlobe,
  'hard-drive': FiHardDrive,
  key: FiKey,
  layers: FiLayers,
  lock: FiLock,
  server: FiServer,
  shield: FiShield,
};

export const topicIcon = (name) => ICONS[name] || FiLayers;
