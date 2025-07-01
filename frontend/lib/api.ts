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
  async getIssues(projectId?: string, limit?: number, offset?: number, filters?: Record<string, string | number | undefined>) {
    let endpoint = projectId ? `/issues?project_id=${projectId}` : '/issues';
    const params = [];
    if (typeof limit === 'number') params.push(`limit=${limit}`);
    if (typeof offset === 'number') params.push(`offset=${offset}`);
    if (filters) {
      for (const [key, value] of Object.entries(filters)) {
        if (value !== undefined && value !== '') params.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`);
      }
    }
    if (params.length > 0) {
      endpoint += (endpoint.includes('?') ? '&' : '?') + params.join('&');
    }
    return this.request<{ issues: any[]; total: number; limit: number; offset: number }>(endpoint)
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
  async getComments(issueId: string, limit?: number, offset?: number) {
    let endpoint = `/comments/issue/${issueId}`;
    const params = [];
    if (typeof limit === 'number') params.push(`limit=${limit}`);
    if (typeof offset === 'number') params.push(`offset=${offset}`);
    if (params.length > 0) {
      endpoint += (endpoint.includes('?') ? '&' : '?') + params.join('&');
    }
    return this.request<{ comments: any[]; total: number; limit: number; offset: number }>(endpoint)
  }

  async createComment(commentData: { content: string; issue_id: string }) {
    return this.request('/comments', {
      method: 'POST',
      body: JSON.stringify(commentData),
    })
  }

  // User endpoints
  async getProfile() {
    return this.request('/users/profile')
  }

  async updateProfile(profileData: { first_name: string; last_name: string; email: string }) {
    return this.request('/users/profile', {
      method: 'PUT',
      body: JSON.stringify(profileData),
    })
  }
}

export const apiClient = new ApiClient() 