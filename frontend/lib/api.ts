const API_BASE_URL = 'http://localhost:8080/api/v1'

export interface ApiResponse<T = any> {
  data?: T
  error?: string
  message?: string
}

class ApiClient {
  private getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem('token')
    return {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    try {
      const url = `${API_BASE_URL}${endpoint}`
      const response = await fetch(url, {
        ...options,
        headers: this.getAuthHeaders(),
      })

      const data = await response.json()

      if (!response.ok) {
        if (response.status === 401) {
          // Handle unauthorized - redirect to login
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          window.location.href = '/login'
          return { error: 'Unauthorized' }
        }
        return { error: data.error || 'Request failed' }
      }

      return { data }
    } catch (error) {
      return { error: 'Network error' }
    }
  }

  // Auth endpoints
  async login(email: string, password: string) {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
  }

  async register(userData: {
    email: string
    password: string
    first_name: string
    last_name: string
  }) {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    })
  }

  // Project endpoints
  async getProjects() {
    return this.request<{ projects: any[] }>('/projects')
  }

  async createProject(projectData: { name: string; description: string }) {
    return this.request('/projects', {
      method: 'POST',
      body: JSON.stringify(projectData),
    })
  }

  async getProject(id: string) {
    return this.request(`/projects/${id}`)
  }

  async updateProject(id: string, projectData: { name: string; description: string }) {
    return this.request(`/projects/${id}`, {
      method: 'PUT',
      body: JSON.stringify(projectData),
    })
  }

  async deleteProject(id: string) {
    return this.request(`/projects/${id}`, {
      method: 'DELETE',
    })
  }

  // Issue endpoints
  async getIssues(projectId?: string) {
    const endpoint = projectId ? `/issues?project_id=${projectId}` : '/issues'
    return this.request<{ issues: any[] }>(endpoint)
  }

  async createIssue(issueData: {
    title: string
    description: string
    project_id: string
    status?: string
    priority?: string
  }) {
    return this.request('/issues', {
      method: 'POST',
      body: JSON.stringify(issueData),
    })
  }

  async getIssue(id: string) {
    return this.request(`/issues/${id}`)
  }

  async updateIssue(id: string, issueData: any) {
    return this.request(`/issues/${id}`, {
      method: 'PUT',
      body: JSON.stringify(issueData),
    })
  }

  async deleteIssue(id: string) {
    return this.request(`/issues/${id}`, {
      method: 'DELETE',
    })
  }

  // Comment endpoints
  async getComments(issueId: string) {
    return this.request<{ comments: any[] }>(`/comments/issue/${issueId}`)
  }

  async createComment(commentData: { content: string; issue_id: string }) {
    return this.request('/comments', {
      method: 'POST',
      body: JSON.stringify(commentData),
    })
  }
}

export const apiClient = new ApiClient() 