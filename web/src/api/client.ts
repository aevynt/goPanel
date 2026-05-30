import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  },
)

export default api

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user_id: number
  username: string
  role: string
}

export interface User {
  id: number
  username: string
  role: string
}

export interface Service {
  name: string
  description: string
  status: string
  port?: number
  binary_path?: string
  panel_managed: boolean
}

export interface ServiceSpec {
  name: string
  description: string
  binary_path: string
  working_dir: string
  port: number
  env_vars: string
  args: string
  auto_start: boolean
  run_as: string
}

export interface Binary {
  id: number
  name: string
  original_name: string
  path: string
  size: number
  version: string
}

export interface FileInfo {
  name: string
  path: string
  size: number
  is_dir: boolean
  mode: string
  mod_time: string
}

export interface PortInfo {
  port: number
  protocol: string
  state: string
  pid?: number
  process?: string
}

export interface Site {
  domain: string
  service_port: number
  tls_enabled: boolean
  tls_email?: string
  extra_config?: string
  type: string
  root?: string
}

export interface DashboardStats {
  uptime: string
  go_version: string
  version: string
  os: string
  kernel: string
  hostname: string
  cpu_percent: number
  memory: {
    total: number
    used: number
    used_percent: number
  }
  disk: {
    total: number
    used: number
    used_percent: number
  }
  load?: {
    load1: number
    load5: number
    load15: number
  }
  services_count: number
  ports_open: number
  sites_count: number
}

export interface PanelSettings {
  panel_domain: string
  port: number
  log_level: string
  public_domain: string
  public_port: number
}

export interface Share {
  id: string
  folder: string
  type: string
  title: string
  description: string
  created_at: string
}

export interface FileEntry {
  name: string
  path: string
  size: number
  is_dir: boolean
  mod_time: string
  thumbnail?: string
  mime_type?: string
}

export interface PublicDomainInfo {
	public_domain: string
	public_port: number
}

export interface UpdateAsset {
	name: string
	content_type: string
	browser_download_url: string
	size: number
}

export interface UpdateRelease {
	tag_name: string
	name: string
	body: string
	html_url: string
	published_at: string
	prerelease: boolean
	assets: UpdateAsset[]
}

export interface UpdateInfo {
	current_version: string
	latest_version: string
	has_update: boolean
	release?: UpdateRelease
	checked_at: string
	error?: string
}
