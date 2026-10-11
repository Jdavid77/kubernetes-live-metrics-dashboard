// Architecture topics shown on /architecture.
// Describe roles, not identifiers: no IPs, ports, bucket names or secret paths.

const topics = [
  {
    id: 'ingress',
    title: 'Traffic Ingress',
    color: 'var(--color-info)',
    summary:
      'No ports are open on the home router. Public apps are reached through a Cloudflare Tunnel, while internal apps only exist in Pi-hole and are reachable from the LAN or over Tailscale.',
    lanes: [
      {
        label: 'Internet',
        steps: ['Client', 'Cloudflare edge', 'cloudflared tunnel', 'Envoy external', 'HTTPRoute', 'App'],
      },
      {
        label: 'LAN',
        steps: ['Client', 'Pi-hole DNS', 'Envoy internal', 'HTTPRoute', 'App'],
      },
      {
        label: 'Remote',
        steps: ['Device', 'Tailscale subnet router', 'Pi-hole DNS', 'Envoy internal', 'HTTPRoute', 'App'],
      },
    ],
    components: [
      { name: 'Cloudflare + cloudflared', role: 'Public DNS, edge TLS/DDoS protection and the outbound tunnel into the cluster' },
      { name: 'Envoy Gateway', role: 'Two Gateway API gateways: external (public) and internal (LAN only)' },
      { name: 'External DNS', role: 'Two instances: external gateway routes go to Cloudflare, internal gateway routes go to Pi-hole' },
      { name: 'Pi-hole', role: 'LAN DNS resolver, also used as the tailnet DNS server' },
      { name: 'Tailscale Operator', role: 'Subnet router advertising the LAN to the tailnet' },
      { name: 'cert-manager', role: "Let's Encrypt certificates for the gateways" },
    ],
    decisions: [
      {
        title: 'Cloudflare Tunnel over port-forwarding',
        why: 'No inbound ports on the router, the router’s public IP stays hidden, and Cloudflare provides TLS and DDoS protection at the edge.',
      },
      {
        title: 'Two gateways instead of one',
        why: 'Exposure is an explicit per-app choice via the route parentRef. An app is LAN-only unless it points at the external gateway, so mistakes fail closed.',
      },
      {
        title: 'Each gateway publishes to its own DNS',
        why: 'Routes on the external gateway are published to Cloudflare, routes on the internal gateway only to Pi-hole. Internal hostnames never appear in public DNS, so they cannot even be discovered from outside.',
      },
      {
        title: 'Tailscale subnet router with Pi-hole as tailnet DNS',
        why: 'Remote devices resolve hostnames exactly as they would at home, so internal apps work remotely with no extra config and nothing exposed publicly.',
      },
    ],
  },
  {
    id: 'secrets',
    title: 'Secrets',
    color: 'var(--color-warning)',
    summary:
      'Secrets never live in Git, not even encrypted. Akeyless is the single source of truth and External Secrets Operator syncs them into the cluster; SOPS is only used for the bootstrap credentials.',
    components: [
      { name: 'Akeyless', role: 'Hosted secrets manager, source of truth for every secret' },
      { name: 'External Secrets Operator', role: 'Pulls secrets from Akeyless into Kubernetes Secrets every 5 minutes' },
      { name: 'SOPS', role: 'Encrypts the few secrets needed before ESO can run' },
      { name: 'Terraform', role: 'Manages Akeyless configuration as code' },
    ],
    decisions: [
      {
        title: 'Akeyless as the source of truth',
        why: 'The free tier covers 500 secrets, plenty for a homelab, and keeps secrets out of Git entirely. Rotation means changing a value in Akeyless and waiting for the next sync.',
      },
      {
        title: 'SOPS only for bootstrap',
        why: 'ESO needs credentials to reach Akeyless before it can fetch anything, so that first secret has to come from somewhere else.',
      },
      {
        title: 'Akeyless configured with Terraform',
        why: 'Auth methods, roles and access rules are reproducible as code instead of clicked together in a UI.',
      },
    ],
  },
  {
    id: 'storage',
    title: 'Storage & Backups',
    color: 'var(--color-success)',
    summary:
      'Every piece of state that matters has a local copy and an offsite copy. Databases live outside the cluster on the NAS, and data that can be rebuilt sits on fast local disks with no backup.',
    components: [
      { name: 'Longhorn', role: 'Replicated block storage, the default for app volumes' },
      { name: 'OpenEBS', role: 'Local hostpath volumes for write-heavy, rebuildable data' },
      { name: 'Volsync', role: 'Snapshot-based PVC backups with Restic' },
      { name: 'Garage', role: 'In-cluster S3 object storage, the local backup target' },
      { name: 'Cloudflare R2 / Backblaze B2', role: 'Offsite object storage for PVC and Postgres backups' },
      { name: 'Postgres', role: 'Runs in Docker on the NAS, outside the cluster' },
    ],
    decisions: [
      {
        title: 'Postgres on the NAS, not in the cluster',
        why: 'Avoids the complexity of a database operator, replicas and failover in-cluster, and frees worker resources for other apps.',
      },
      {
        title: 'Local + offsite copy of everything',
        why: 'Local copies make restores fast; offsite copies survive losing the whole house.',
      },
      {
        title: 'Longhorn by default, OpenEBS for rebuildable data',
        why: 'Metrics, logs, scan reports are write-heavy and cheap to regenerate, so they skip Longhorn network replication and use local disks.',
      },
    ],
  },
  {
    id: 'observability',
    title: 'Observability',
    color: 'var(--color-secondary)',
    summary:
      'Metrics, logs and uptime checks feed Grafana, and every alert ends up in a single notification channel. Monitoring is declared next to the thing it monitors.',
    components: [
      { name: 'kube-prometheus-stack', role: 'Prometheus, Alertmanager and the core exporters' },
      { name: 'Grafana Operator', role: 'Grafana with dashboards defined as CRs' },
      { name: 'OpenTelemetry Collector', role: 'Collects pod logs, Kubernetes events and GitHub data' },
      { name: 'Loki', role: 'Log storage and querying' },
      { name: 'Gatus', role: 'Uptime checks declared as annotations on gateways and apps' },
      { name: 'Pushover', role: 'Push notifications to phone' },
    ],
    decisions: [
      {
        title: 'Pushover as the single notification channel',
        why: 'Prometheus rules and Flux reconcile failures route through Alertmanager, and Gatus notifies directly, so every alert lands in one place.',
      },
      {
        title: 'Monitoring declared alongside each app',
        why: 'Gatus annotations, Grafana dashboards and alert components live with the app manifests, so adding an app brings its monitoring with it.',
      },
      {
        title: 'OpenTelemetry over Promtail',
        why: 'Promtail is deprecated and OpenTelemetry is the vendor-neutral industry standard',
      },
    ],
  },
];

export default topics;
