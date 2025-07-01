"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Link from "next/link";
import { apiClient } from "../../../../lib/api";

export default function EditProjectPage() {
  const router = useRouter();
  const params = useParams();
  const projectId = typeof params?.id === "string" ? params.id : Array.isArray(params?.id) ? params.id[0] : "";

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    if (!projectId) return;
    fetchProject();
    // eslint-disable-next-line
  }, [projectId]);

  const fetchProject = async () => {
    setIsLoading(true);
    setError("");
    try {
      const response = await apiClient.getProject(projectId);
      if (response.error || !response.data) {
        setError(response.error || "Project not found");
      } else {
        const project = response.data as { name: string; description: string };
        setName(project.name);
        setDescription(project.description);
      }
    } catch (err) {
      setError("Failed to fetch project");
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess(false);
    try {
      const response = await apiClient.updateProject(projectId, { name, description });
      if (response.error) {
        setError(response.error);
      } else {
        setSuccess(true);
        setTimeout(() => {
          router.push(`/projects/${projectId}`);
        }, 1000);
      }
    } catch (err) {
      setError("Failed to update project");
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md mx-auto">
        <Link href={`/projects/${projectId}`} className="text-blue-600 hover:text-blue-800 font-medium mb-4 inline-block">
          ← Back to Project
        </Link>
        <div className="card p-6">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">Edit Project</h1>
          {error && <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4">{error}</div>}
          {success && <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-md mb-4">Project updated!</div>}
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700">Project Name</label>
              <input
                id="name"
                type="text"
                value={name}
                onChange={e => setName(e.target.value)}
                className="input-field mt-1"
                required
              />
            </div>
            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700">Description</label>
              <textarea
                id="description"
                value={description}
                onChange={e => setDescription(e.target.value)}
                className="input-field mt-1"
                rows={4}
              />
            </div>
            <div>
              <button type="submit" className="btn-primary w-full">Save Changes</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
} 