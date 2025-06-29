'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { apiClient } from '../../../lib/api'

interface Project {
  id: string
  name: string
  description: string
  created_at: string
  updated_at: string
}

interface Issue {
  id: string
  title: string
  description: string
  status: string
  priority: string
  project_id: string
  created_at: string
  updated_at: string
}

export default function ProjectDetailPage({ params }: { params: { id: string } }) {
  const router = useRouter()
  const [project, setProject] = useState<Project | null>(null)
  const [issues, setIssues] = useState<Issue[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    fetchProject()
    fetchIssues()
  }, [params.id])

  const fetchProject = async () => {
    try {
      const response = await apiClient.getProject(params.id)
      
      if (response.error) {
        setError(response.error)
      } else {
        setProject(response.data as Project)
      }
    } catch (err) {
      setError('Failed to fetch project')
    }
  }

  const fetchIssues = async () => {
    try {
      const response = await apiClient.getIssues(params.id)
      
      if (response.error) {
        console.error('Failed to fetch issues:', response.error)
      } else {
        setIssues(response.data?.issues || [])
      }
    } catch (err) {
      console.error('Failed to fetch issues:', err)
    } finally {
      setIsLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete this project? This action cannot be undone.')) {
      return
    }

    try {
      const response = await apiClient.deleteProject(params.id)
      
      if (response.error) {
        setError(response.error)
      } else {
        router.push('/dashboard')
      }
    } catch (err) {
      setError('Failed to delete project')
    }
  }

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'critical': return 'bg-red-100 text-red-800'
      case 'high': return 'bg-orange-100 text-orange-800'
      case 'medium': return 'bg-yellow-100 text-yellow-800'
      case 'low': return 'bg-green-100 text-green-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'open': return 'bg-blue-100 text-blue-800'
      case 'in_progress': return 'bg-yellow-100 text-yellow-800'
      case 'resolved': return 'bg-green-100 text-green-800'
      case 'closed': return 'bg-gray-100 text-gray-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-2xl mx-auto">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
            {error}
          </div>
          <Link 
            href="/dashboard" 
            className="text-blue-600 hover:text-blue-800 font-medium mt-4 inline-block"
          >
            ← Back to Dashboard
          </Link>
        </div>
      </div>
    )
  }

  if (!project) {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-2xl mx-auto">
          <div className="text-center">
            <h1 className="text-2xl font-bold text-gray-900 mb-4">Project not found</h1>
            <Link 
              href="/dashboard" 
              className="text-blue-600 hover:text-blue-800 font-medium"
            >
              ← Back to Dashboard
            </Link>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <Link 
            href="/dashboard" 
            className="text-blue-600 hover:text-blue-800 font-medium mb-4 inline-block"
          >
            ← Back to Dashboard
          </Link>
          <div className="flex justify-between items-start">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">{project.name}</h1>
              <p className="mt-2 text-gray-600">
                Created on {new Date(project.created_at).toLocaleDateString()}
              </p>
            </div>
            <div className="flex space-x-4">
              <Link
                href={`/projects/${project.id}/edit`}
                className="btn-primary"
              >
                Edit Project
              </Link>
              <button
                onClick={handleDelete}
                className="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded-md transition-colors duration-200"
              >
                Delete Project
              </button>
            </div>
          </div>
        </div>

        {/* Project details */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main content */}
          <div className="lg:col-span-2">
            <div className="card p-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">Description</h2>
              <p className="text-gray-700 whitespace-pre-wrap">
                {project.description || 'No description provided.'}
              </p>
            </div>

            {/* Issues section */}
            <div className="card p-6 mt-6">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl font-semibold text-gray-900">Issues ({issues.length})</h2>
                <Link
                  href={`/projects/${project.id}/issues/new`}
                  className="btn-primary"
                >
                  Create Issue
                </Link>
              </div>
              
              {issues.length === 0 ? (
                <div className="text-center py-8 text-gray-500">
                  <p>No issues yet. Create the first issue to get started!</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {issues.map((issue) => (
                    <div key={issue.id} className="border border-gray-200 rounded-lg p-4 hover:bg-gray-50">
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <h3 className="text-lg font-medium text-gray-900 mb-2">
                            <Link 
                              href={`/issues/${issue.id}`}
                              className="hover:text-blue-600"
                            >
                              {issue.title}
                            </Link>
                          </h3>
                          <p className="text-gray-600 text-sm line-clamp-2">
                            {issue.description}
                          </p>
                        </div>
                        <div className="flex space-x-2 ml-4">
                          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getPriorityColor(issue.priority)}`}>
                            {issue.priority}
                          </span>
                          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(issue.status)}`}>
                            {issue.status.replace('_', ' ')}
                          </span>
                        </div>
                      </div>
                      <div className="mt-3 text-sm text-gray-500">
                        Created {new Date(issue.created_at).toLocaleDateString()}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="lg:col-span-1">
            <div className="card p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Project Info</h3>
              <dl className="space-y-3">
                <div>
                  <dt className="text-sm font-medium text-gray-500">Created</dt>
                  <dd className="text-sm text-gray-900">
                    {new Date(project.created_at).toLocaleDateString()}
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Last Updated</dt>
                  <dd className="text-sm text-gray-900">
                    {new Date(project.updated_at).toLocaleDateString()}
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Total Issues</dt>
                  <dd className="text-sm text-gray-900">
                    {issues.length}
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Status</dt>
                  <dd className="text-sm text-gray-900">
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      Active
                    </span>
                  </dd>
                </div>
              </dl>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
} 