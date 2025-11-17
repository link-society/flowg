export const ScopeLabels = {
  read_pipelines: 'View Pipelines',
  write_pipelines: 'View & Edit Pipelines',
  read_transformers: 'View Transformers',
  write_transformers: 'View & Edit Transformers',
  read_streams: 'View Streams',
  write_streams: 'View & Edit Streams',
  read_forwarders: 'View Forwarders',
  write_forwarders: 'View & Edit Forwarders',
  read_acls: 'View ACLs',
  write_acls: 'View & Edit ACLs',
  send_logs: 'Send Logs',
}

export const Scopes = Object.keys(ScopeLabels) as string[]
