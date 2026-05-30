import api, {
  type LoginRequest,
  type LoginResponse,
  type User,
  type Service,
  type ServiceSpec,
  type Binary,
  type FileInfo,
  type PortInfo,
  type Site,
  type DashboardStats,
  type Share,
  type FileEntry,
  type PublicDomainInfo,
} from './client'
import type { PanelSettings, UpdateInfo } from './client'

export async function login(data: LoginRequest): Promise<LoginResponse> {
  const res = await api.post('/auth/login', data)
  return res.data
}

export async function refreshToken(token: string): Promise<string> {
  const res = await api.post('/auth/refresh', { token })
  return res.data.token
}

export async function getMe(): Promise<LoginResponse> {
  const res = await api.get('/auth/me')
  return res.data
}

export async function listUsers(): Promise<User[]> {
  const res = await api.get('/users')
  return res.data ?? []
}

export async function createUser(data: { username: string; password: string; role: string }): Promise<User> {
  const res = await api.post('/users', data)
  return res.data
}

export async function updateUser(id: number, data: { username?: string; password?: string; role?: string }): Promise<void> {
  await api.put(`/users/${id}`, data)
}

export async function deleteUser(id: number): Promise<void> {
  await api.delete(`/users/${id}`)
}

export async function listServices(): Promise<Service[]> {
  const res = await api.get('/services')
  return res.data
}

export async function getService(name: string): Promise<Service> {
  const res = await api.get(`/services/${name}`)
  return res.data
}

export async function createService(spec: ServiceSpec): Promise<void> {
  await api.post('/services', spec)
}

export async function startService(name: string): Promise<void> {
  await api.post(`/services/${name}/start`)
}

export async function stopService(name: string): Promise<void> {
  await api.post(`/services/${name}/stop`)
}

export async function restartService(name: string): Promise<void> {
  await api.post(`/services/${name}/restart`)
}

export async function enableService(name: string): Promise<void> {
  await api.post(`/services/${name}/enable`)
}

export async function disableService(name: string): Promise<void> {
  await api.post(`/services/${name}/disable`)
}

export async function removeService(name: string): Promise<void> {
  await api.delete(`/services/${name}`)
}

export async function getServiceLogs(name: string, tail = 50): Promise<{ timestamp: string; message: string }[]> {
  const res = await api.get(`/services/${name}/logs`, { params: { tail } })
  return res.data
}

export async function listBinaries(): Promise<Binary[]> {
  const res = await api.get('/binaries')
  return res.data ?? []
}

export async function uploadBinary(file: File, version?: string, onProgress?: (pct: number) => void): Promise<Binary> {
  const form = new FormData()
  form.append('file', file)
  if (version) form.append('version', version)
  const res = await api.post('/binaries', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: (e) => {
      if (e.total && onProgress) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    },
  })
  return res.data
}

export async function deleteBinary(id: number): Promise<void> {
  await api.delete(`/binaries/${id}`)
}

export async function listFiles(path = '/'): Promise<FileInfo[]> {
  const res = await api.get('/files', { params: { path } })
  return res.data ?? []
}

export async function readFile(path: string): Promise<string> {
  const res = await api.get('/files/read', { params: { path }, responseType: 'text' })
  return res.data
}

export async function getFileBlob(path: string): Promise<Blob> {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api/v1/files/read?path=${encodeURIComponent(path)}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) throw new Error('Failed to load file')
  return res.blob()
}

export async function writeFile(path: string, content: string): Promise<void> {
  await api.post('/files/write', { path, content })
}

export async function mkdir(path: string): Promise<void> {
  await api.post('/files/mkdir', { path })
}

export async function uploadFile(path: string, file: File, onProgress?: (pct: number) => void): Promise<void> {
  const form = new FormData()
  form.append('file', file)
  form.append('path', path)
  await api.post('/files/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: (e) => {
      if (e.total && onProgress) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    },
  })
}

export async function removeFile(path: string): Promise<void> {
  await api.delete('/files', { data: { path } })
}

export async function renameFile(oldPath: string, newPath: string): Promise<void> {
  await api.post('/files/rename', { old_path: oldPath, new_path: newPath })
}

export async function listPorts(start?: number, end?: number): Promise<PortInfo[]> {
  const res = await api.get('/ports', { params: { start, end } })
  return res.data
}

export async function checkPort(port: number): Promise<{ port: number; available: boolean }> {
  const res = await api.get(`/ports/check/${port}`)
  return res.data
}

export async function findPort(preferred?: number): Promise<{ port: number }> {
  const res = await api.post('/ports/find', { preferred })
  return res.data
}

export async function listSites(): Promise<Site[]> {
  const res = await api.get('/sites')
  return res.data
}

export async function addSite(site: Site): Promise<void> {
  await api.post('/sites', site)
}

export async function removeSite(domain: string): Promise<void> {
  await api.delete(`/sites/${domain}`)
}

export async function caddyHealth(): Promise<void> {
  await api.get('/sites/health')
}

export async function getDashboardStats(): Promise<DashboardStats> {
  const res = await api.get('/dashboard/stats')
  return res.data
}

export async function getSettings(): Promise<PanelSettings> {
  const res = await api.get('/settings')
  return res.data
}

export async function updateSettings(data: Partial<PanelSettings>): Promise<PanelSettings> {
  const res = await api.put('/settings', data)
  return res.data
}

export async function listShares(): Promise<Share[]> {
  const res = await api.get('/public/shares')
  return res.data ?? []
}

export async function createShare(data: { id?: string; title?: string; description?: string }): Promise<Share> {
  const res = await api.post('/public/shares', data)
  return res.data
}

export async function deleteShare(id: string): Promise<void> {
  await api.delete(`/public/shares/${id}`)
}

export async function listShareFiles(id: string): Promise<FileEntry[]> {
  const res = await api.get(`/public/shares/${id}/files`)
  return res.data ?? []
}

export async function uploadShareFile(id: string, files: File[]): Promise<void> {
  const form = new FormData()
  files.forEach(f => form.append('files', f))
  await api.post(`/public/shares/${id}/upload`, form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export async function deleteShareFile(id: string, name: string): Promise<void> {
  await api.delete(`/public/shares/${id}/file`, { data: { id, name } })
}

export async function getPublicDomain(): Promise<PublicDomainInfo> {
  const res = await api.get('/public/domain')
  return res.data
}

export async function setPublicDomain(domain: string): Promise<void> {
	await api.put('/public/domain', { public_domain: domain })
}

export async function checkUpdate(): Promise<UpdateInfo> {
	const res = await api.get('/updates/check')
	return res.data
}
