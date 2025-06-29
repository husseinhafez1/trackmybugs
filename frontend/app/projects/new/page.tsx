'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { apiClient } from '../../../lib/api'

interface ProjectForm {
  name: string
  description: string
}

export default function NewProjectPage() {
  const router = useRouter()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  
  const { register, handleSubmit, formState: { errors } } = useForm<ProjectForm>()

  const onSubmit = async (data: ProjectForm) => {
    setIsLoading(true)
    setError('')

    try {
      const response = await apiClient.createProject(data)
      
      if (response.error) {
        setError(response.error)
      } else {
        router.push('/dashboard')
      }
    } catch (err) {
      setError('Failed to create project')
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
            href="/dashboard" 
            className="text-blue-600 hover:text-blue-800 font-medium mb-4 inline-block"
          >
            ← Back to Dashboard
          </Link>
          <h1 className="text-3xl font-bold text-gray-900">Create New Project</h1>
          <p className="mt-2 text-gray-600">
            Set up a new project to start tracking bugs and issues.
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
              <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                Project Name *
              </label>
              <input
                {...register('name', { 
                  required: 'Project name is required',
                  minLength: {
                    value: 2,
                    message: 'Project name must be at least 2 characters'
                  }
                })}
                type="text"
                className="input-field mt-1"
                placeholder="Enter project name"
              />
              {errors.name && (
                <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                Description
              </label>
              <textarea
                {...register('description', {
                  maxLength: {
                    value: 500,
                    message: 'Description must be less than 500 characters'
                  }
                })}
                rows={4}
                className="input-field mt-1"
                placeholder="Describe your project (optional)"
              />
              {errors.description && (
                <p className="mt-1 text-sm text-red-600">{errors.description.message}</p>
              )}
            </div>

            <div className="flex justify-end space-x-4">
              <Link
                href="/dashboard"
                className="btn-secondary"
              >
                Cancel
              </Link>
              <button
                type="submit"
                disabled={isLoading}
                className="btn-primary"
              >
                {isLoading ? 'Creating...' : 'Create Project'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
} 