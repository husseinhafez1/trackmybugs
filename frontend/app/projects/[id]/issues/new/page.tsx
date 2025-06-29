'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { apiClient } from '../../../../../lib/api'

interface IssueForm {
  title: string
  description: string
  priority: string
  status: string
}

export default function NewIssuePage({ params }: { params: { id: string } }) {
  const router = useRouter()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  
  const { register, handleSubmit, formState: { errors } } = useForm<IssueForm>({
    defaultValues: {
      priority: 'medium',
      status: 'open'
    }
  })

  const onSubmit = async (data: IssueForm) => {
    setIsLoading(true)
    setError('')

    try {
      const response = await apiClient.createIssue({
        ...data,
        project_id: params.id
      })
      
      if (response.error) {
        setError(response.error)
      } else {
        router.push(`/projects/${params.id}`)
      }
    } catch (err) {
      setError('Failed to create issue')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <Link 
            href={`/projects/${params.id}`}
            className="text-blue-600 hover:text-blue-800 font-medium mb-4 inline-block"
          >
            ← Back to Project
          </Link>
          <h1 className="text-3xl font-bold text-gray-900">Create New Issue</h1>
          <p className="mt-2 text-gray-600">
            Report a new bug or feature request for this project.
          </p>
        </div>

        {/* Form */}
        <div className="card p-6">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
                {error}
              </div>
            )}
            
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700">
                Issue Title *
              </label>
              <input
                {...register('title', { 
                  required: 'Issue title is required',
                  minLength: {
                    value: 3,
                    message: 'Title must be at least 3 characters'
                  }
                })}
                type="text"
                className="input-field mt-1"
                placeholder="Brief description of the issue"
              />
              {errors.title && (
                <p className="mt-1 text-sm text-red-600">{errors.title.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                Description *
              </label>
              <textarea
                {...register('description', { 
                  required: 'Description is required',
                  minLength: {
                    value: 10,
                    message: 'Description must be at least 10 characters'
                  }
                })}
                rows={6}
                className="input-field mt-1"
                placeholder="Detailed description of the issue, steps to reproduce, expected vs actual behavior..."
              />
              {errors.description && (
                <p className="mt-1 text-sm text-red-600">{errors.description.message}</p>
              )}
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="priority" className="block text-sm font-medium text-gray-700">
                  Priority
                </label>
                <select
                  {...register('priority')}
                  className="input-field mt-1"
                >
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                  <option value="critical">Critical</option>
                </select>
              </div>

              <div>
                <label htmlFor="status" className="block text-sm font-medium text-gray-700">
                  Status
                </label>
                <select
                  {...register('status')}
                  className="input-field mt-1"
                >
                  <option value="open">Open</option>
                  <option value="in_progress">In Progress</option>
                  <option value="resolved">Resolved</option>
                  <option value="closed">Closed</option>
                </select>
              </div>
            </div>

            <div className="flex justify-end space-x-4">
              <Link
                href={`/projects/${params.id}`}
                className="btn-secondary"
              >
                Cancel
              </Link>
              <button
                type="submit"
                disabled={isLoading}
                className="btn-primary"
              >
                {isLoading ? 'Creating...' : 'Create Issue'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
} 