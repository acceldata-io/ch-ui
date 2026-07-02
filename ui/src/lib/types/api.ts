/** Standard API response envelope */
export interface ApiResponse<T = unknown> {
  success: boolean
  error?: string
  data?: T
}

/** Session info returned by /api/auth/session */
export interface Session {
  user: string
  role: string
  connectionId: string
  connectionName: string
  connectionOnline: boolean
  expiresAt: string
  version?: string
  appVersion?: string
}

/** Connection info */
export interface Connection {
  id: string
  name: string
  status: string
  online: boolean
  created_at: string
  host_info?: HostInfo
}

/** Host machine metrics from agent */
export interface HostInfo {
  hostname: string
  os: string
  arch: string
  cpu_cores: number
  memory_total: number
  memory_free: number
  disk_total: number
  disk_free: number
  go_version: string
  agent_uptime: number
  collected_at: string
}

/** Saved query */
export interface SavedQuery {
  id: string
  name: string
  query: string
  description?: string
  parameters?: string | null // JSON object of default {name: value} bind params
  created_at: string
  updated_at: string
}

/** Dashboard */
export interface Dashboard {
  id: string
  name: string
  description: string | null
  created_by: string
  created_at: string
  updated_at: string
}

/** Dashboard panel */
export interface Panel {
  id: string
  dashboard_id: string
  name: string
  description: string
  panel_type: string
  query: string
  connection_id: string | null
  config: string
  layout_x: number
  layout_y: number
  layout_w: number
  layout_h: number
  created_at: string
  updated_at: string
}

export interface StatThreshold {
  value: number
  color: string
}

/** Panel visualization config (stored as JSON in panel.config) */
export interface PanelConfig {
  chartType: 'table' | 'stat' | 'timeseries' | 'bar' | 'text' | 'gauge' | 'pie'
  xColumn?: string
  yColumns?: string[]
  colors?: string[]
  legendPosition?: 'bottom' | 'right' | 'none'
  content?: string
  barMode?: 'grouped' | 'stacked'
  gaugeMin?: number
  gaugeMax?: number
  pieDonut?: boolean
  pieLabelColumn?: string
  pieValueColumn?: string
  statField?: string
  statCalculation?: 'last' | 'first' | 'mean' | 'sum' | 'min' | 'max' | 'count' | 'range'
  statUnit?: 'none' | 'percent' | 'short' | 'bytes' | 'bps' | 'duration' | 'durationMs'
  statPrefix?: string
  statSuffix?: string
  statDecimals?: number
  statColorMode?: 'none' | 'value' | 'background'
  statThresholds?: StatThreshold[]
}

/** Audit log entry */
export interface AuditLog {
  id: string
  action: string
  username: string | null
  details: string | null
  ip_address: string | null
  created_at: string
  parsed_details?: Record<string, unknown>
}

/** Admin stats overview */
export interface AdminStats {
  users_count: number
  connections: number
  online: number
  login_count: number
  query_count: number
}
