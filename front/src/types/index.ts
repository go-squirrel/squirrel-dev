// 全局类型定义

// IP地址类型
export interface IpAddress {
  ipv4: string[]
  ipv6: string[]
  name: string
}

// 服务器详细信息
export interface ServerInfo {
  architecture: string
  hostname: string
  ipAddresses: IpAddress[]
  kernelVersion: string
  os: string
  platform: string
  platformVersion: string
  uptime: number
  uptimeStr: string
}

// 服务器类型
export interface Server {
  id: number
  hostname: string
  ip_address: string
  ssh_username: string
  ssh_port: number
  auth_type: 'password' | 'key'
  status: 'online' | 'offline' | 'unknown' | 'active' | 'inactive'
  server_info?: ServerInfo | null
  server_alias?: string
}

// 创建服务器请求
export interface CreateServerRequest {
  ip_address: string
  ssh_username: string
  ssh_port: number
  ssh_password?: string
  ssh_private_key?: string
  auth_type: 'password' | 'key'
  status: 'active' | 'inactive'
  server_alias?: string
}

// 更新服务器请求
export interface UpdateServerRequest {
  ip_address: string
  ssh_username: string
  ssh_port: number
  auth_type: 'password' | 'key'
  status: 'active' | 'inactive'
  server_alias?: string
  ssh_password?: string
  ssh_private_key?: string
}

// 应用类型
export interface Application {
  id: number
  name: string
  version: string
  status: 'running' | 'stopped' | 'error'
  icon: string
  color: string
}

// 监控数据类型
export interface MonitorData {
  hostname?: string
  cpu?: {
    cores: number
    frequency: number
    model: string
    perCoreUsage: number[]
    usage: number
  }
  memory?: {
    available: number
    swapTotal: number
    swapUsed: number
    total: number
    usage: number
    used: number
  }
  disk?: {
    available: number
    partitions: DiskPartition[]
    total: number
    usage: number
    used: number
  }
  loadAverage?: {
    load1: number
    load5: number
    load15: number
  }
  topCPU?: ProcessInfo[]
  topMemory?: ProcessInfo[]
}

export interface DiskPartition {
  available: number
  device: string
  fsType: string
  mountPoint: string
  total: number
  usage: number
  used: number
}

export interface ProcessInfo {
  pid: number
  name: string
  cpuPercent: number
  memoryMB: number
  memoryPercent: number
  status: string
  createTime: number
}

// 系统信息类型
export interface SystemInfo {
  hostname: string
  os: string
  kernel: string
  arch: string
  ip: string
  bootTime: string
  uptime: string
}

// 用户类型
export interface User {
  id: number
  username: string
  avatar?: string
  role: string
}

// API 响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 图表数据点
export interface ChartDataPoint {
  time: string
  value1: number
  value2: number
}

// 配置项类型
export interface Config {
  id: number
  key: string
  value: string
}

// 创建配置请求
export interface CreateConfigRequest {
  key: string
  value: string
}

// 更新配置请求
export interface UpdateConfigRequest {
  id: number
  key: string
  value: string
}

// 脚本类型
export interface Script {
  id: number
  name: string
  content: string
}

// 创建脚本请求
export interface CreateScriptRequest {
  name: string
  content: string
}

// 更新脚本请求
export interface UpdateScriptRequest {
  id: number
  name: string
  content: string
}

// 执行脚本请求
export interface ExecuteScriptRequest {
  script_id: number
  server_id: number
}

// 脚本执行结果
export interface ScriptResult {
  id: number
  task_id: number
  script_id: number
  server_id: number
  server_ip: string
  agent_port: number
  output: string
  status: 'running' | 'success' | 'failed'
  error_message: string
  created_at: string
}

// 应用模板类型
export interface AppStore {
  id: number
  name: string
  description: string
  type: 'compose' | 'k8s_manifest' | 'helm_chart'
  category: 'web' | 'database' | 'middleware' | 'devops'
  icon?: string
  version: string
  content: string
  tags: string
  author: string
  repo_url?: string
  homepage_url?: string
  is_official: boolean
  downloads: number
  status: 'active' | 'deprecated'
}

// 创建应用请求
export interface CreateAppRequest {
  name: string
  description: string
  type: 'compose' | 'k8s_manifest' | 'helm_chart'
  category: 'web' | 'database' | 'middleware' | 'devops'
  icon?: string
  version: string
  content: string
  tags: string
  author: string
  repo_url?: string
  homepage_url?: string
  is_official: boolean
}

// 更新应用请求
export interface UpdateAppRequest {
  id: number
  name: string
  description: string
  type: 'compose' | 'k8s_manifest' | 'helm_chart'
  category: 'web' | 'database' | 'middleware' | 'devops'
  icon?: string
  version: string
  content: string
  tags: string
  author: string
  repo_url?: string
  homepage_url?: string
  is_official: boolean
}
