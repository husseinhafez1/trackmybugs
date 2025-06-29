'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { apiClient } from '../../../lib/api'

interface Issue {
  id: string
  title: string
  description: string
  status: string
  priority: string
  project_id: string
  created_by: string
  created_at: string
  updated_at: string
}

export default function IssueDetailPage({ params }: { params: { id: string } }) {
  const router = useRouter()
  const [issue, setIssue] = useState<Issue | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [isEditing, setIsEditing] = useState(false)
  const [editForm, setEditForm] = useState({
    title: '',
    description: '',
    status: '',
    priority: ''
  })

  useEffect(() => {
    fetchIssue()
  }, [params.id])

  const fetchIssue = async () => {
    try {
      const response = await apiClient.getIssue(params.id)
      
      if (response.error) {
        setError(response.error)
      } else {
        const issueData = response.data as Issue
        setIssue(issueData)
        setEditForm({
          title: issueData.title,
          description: issueData.description,
          status: issueData.status,
          priority: issueData.priority
        })
      }
    } catch (err) {
      setError('Failed to fetch issue')
    } finally {
      setIsLoading(false)
    }
  }

  const handleUpdate = async () => {
    if (!issue) return

    try {
      const response = await apiClient.updateIssue(issue.id, editForm)
      
      if (response.error) {
        setError(response.error)
      } else {
        setIssue({ ...issue, ...editForm })
        setIsEditing(false)
      }
    } catch (err) {
      setError('Failed to update issue')
    }
  }

  const handleDelete = async () => {
    if (!issue || !confirm('Are you sure you want to delete this issue? This action cannot be undone.')) {
      return
    }

    try {
      const response = await apiClient.deleteIssue(issue.id)
      
      if (response.error) {
        setError(response.error)
      } else {
        router.push(`/projects/${issue.project_id}`)
      }
    } catch (err) {
      setError('Failed to delete issue')
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

  if (!issue) {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-2xl mx-auto">
          <div className="text-center">
            <h1 className="text-2xl font-bold text-gray-900 mb-4">Issue not found</h1>
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
            href={`/projects/${issue.project_id}`}
            className="text-blue-600 hover:text-blue-800 font-medium mb-4 inline-block"
          >
            ← Back to Project
          </Link>
          <div className="flex justify-between items-start">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                {isEditing ? (
                  <input
                    type="text"
                    value={editForm.title}
                    onChange={(e) => setEditForm({ ...editForm, title: e.target.value })}
                    className="input-field text-3xl font-bold"
                  />
                ) : (
                  issue.title
                )}
              </h1>
              <p className="mt-2 text-gray-600">
                Created on {new Date(issue.created_at).toLocaleDateString()}
              </p>
            </div>
            <div className="flex space-x-4">
              {isEditing ? (
                <>
                  <button
                    onClick={handleUpdate}
                    className="btn-primary"
                  >
                    Save Changes
                  </button>
                  <button
                    onClick={() => {
                      setIsEditing(false)
                      setEditForm({
                        title: issue.title,
                        description: issue.description,
                        status: issue.status,
                        priority: issue.priority
                      })
                    }}
                    className="btn-secondary"
                  >
                    Cancel
                  </button>
                </>
              ) : (
                <>
                  <button
                    onClick={() => setIsEditing(true)}
                    className="btn-primary"
                  >
                    Edit Issue
                  </button>
                  <button
                    onClick={handleDelete}
                    className="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded-md transition-colors duration-200"
                  >
                    Delete Issue
                  </button>
                </>
              )}
            </div>
          </div>
        </div>

        {/* Issue details */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main content */}
          <div className="lg:col-span-2">
            <div className="card p-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">Description</h2>
              {isEditing ? (
                <textarea
                  value={editForm.description}
                  onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
                  rows={8}
                  className="input-field"
                />
              ) : (
                <p className="text-gray-700 whitespace-pre-wrap">
                  {issue.description}
                </p>
              )}
            </div>

            {/* Comments section - placeholder for now */}
            <div className="card p-6 mt-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">Comments</h2>
              <div className="text-center py-8 text-gray-500">
                <p>Comments feature coming soon!</p>
              </div>
            </div>
          </div>

          {/* Sidebar */}
          <div className="lg:col-span-1">
            <div className="card p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Issue Info</h3>
              <dl className="space-y-3">
                <div>
                  <dt className="text-sm font-medium text-gray-500">Status</dt>
                  <dd className="text-sm text-gray-900">
                    {isEditing ? (
                      <select
                        value={editForm.status}
                        onChange={(e) => setEditForm({ ...editForm, status: e.target.value })}
                        className="input-field"
                      >
                        <option value="open">Open</option>
                        <option value="in_progress">In Progress</option>
                        <option value="resolved">Resolved</option>
                        <option value="closed">Closed</option>
                      </select>
                    ) : (
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(issue.status)}`}>
                        {issue.status.replace('_', ' ')}
                      </span>
                    )}
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Priority</dt>
                  <dd className="text-sm text-gray-900">
                    {isEditing ? (
                      <select
                        value={editForm.priority}
                        onChange={(e) => setEditForm({ ...editForm, priority: e.target.value })}
                        className="input-field"
                      >
                        <option value="low">Low</option>
                        <option value="medium">Medium</option>
                        <option value="high">High</option>
                        <option value="critical">Critical</option>
                      </select>
                    ) : (
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getPriorityColor(issue.priority)}`}>
                        {issue.priority}
                      </span>
                    )}
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Created</dt>
                  <dd className="text-sm text-gray-900">
                    {new Date(issue.created_at).toLocaleDateString()}
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Last Updated</dt>
                  <dd className="text-sm text-gray-900">
                    {new Date(issue.updated_at).toLocaleDateString()}
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