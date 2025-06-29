// frontend/app/dashboard/page.tsx
'use client';

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'

interface Project {
  id: string
  name: string
  description: string
  created_at: string
  updated_at: string
}

interface User {
  id: string
  email: string
  first_name: string
  last_name: string
}

export default function DashboardPage() {
  const router = useRouter()
  const [projects, setProjects] = useState<Project[]>([])
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const token = localStorage.getItem('token')
    const userData = localStorage.getItem('user')

    if (!token) {
      router.push('/login')
      return
    }

    if (userData) {
      setUser(JSON.parse(userData))
    }

    fetchProjects(token)
  }, [router])

  const fetchProjects = async (token: string) => {
    try {
      const response = await fetch('http://localhost:8080/api/v1/projects', {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      if (response.status === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        router.push('/login')
        return
      }

      const data = await response.json()
      if (response.ok) {
        setProjects(data.projects || [])
      } else {
        setError(data.error || 'Failed to fetch projects')
      }
    } catch (err) {
      setError('Network error. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/login')
  }

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div className="flex items-center">
              <h1 className="text-3xl font-bold text-gray-900">TrackMyBugs</h1>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-gray-700">
                Welcome, {user?.first_name} {user?.last_name}
              </span>
              <button
                onClick={handleLogout}
                className="btn-secondary"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main content */}
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {error && (
          <div className="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
            {error}
          </div>
        )}

        {/* Quick actions */}
        <div className="mb-8">
          <div className="flex justify-between items-center">
            <h2 className="text-2xl font-bold text-gray-900">Projects</h2>
            <Link href="/projects/new" className="btn-primary">
              Create New Project
            </Link>
          </div>
        </div>

        {/* Projects grid */}
        {projects.length === 0 ? (
          <div className="text-center py-12">
            <div className="text-gray-500 mb-4">
              <svg className="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">No projects yet</h3>
            <p className="text-gray-500 mb-4">Get started by creating your first project.</p>
            <Link href="/projects/new" className="btn-primary">
              Create Project
            </Link>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {projects.map((project) => (
              <div key={project.id} className="card p-6">
                <div className="flex justify-between items-start mb-4">
                  <h3 className="text-lg font-semibold text-gray-900">{project.name}</h3>
                </div>
                <p className="text-gray-600 mb-4 line-clamp-3">{project.description}</p>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-500">
                    Updated {new Date(project.updated_at).toLocaleDateString()}
                  </span>
                  <Link
                    href={`/projects/${project.id}`}
                    className="text-blue-600 hover:text-blue-800 font-medium"
                  >
                    View Project →
                  </Link>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}