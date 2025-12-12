import { useState } from 'react'
import type { CreateAttendeeRequest } from '../api/types'

interface AttendeeFormProps {
  onSubmit: (data: CreateAttendeeRequest) => Promise<void>
  onCancel: () => void
}

export function AttendeeForm({ onSubmit, onCancel }: AttendeeFormProps) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    const form = e.currentTarget
    const formData = new FormData(form)

    const data: CreateAttendeeRequest = {
      name: formData.get('name') as string,
      email: formData.get('email') as string,
      status: formData.get('status') as 'attending' | 'maybe' | 'declined',
    }

    try {
      await onSubmit(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add attendee')
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="bg-cream-200 rounded-lg p-4 border border-farm-200">
      {error && (
        <div className="mb-4 p-3 bg-rust-100 text-rust-800 rounded text-sm border border-rust-300">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label className="block text-sm font-medium text-farm-700 mb-1">
            Name *
          </label>
          <input
            type="text"
            name="name"
            required
            className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 focus:border-transparent bg-cream-50"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-farm-700 mb-1">
            Email *
          </label>
          <input
            type="email"
            name="email"
            required
            className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 focus:border-transparent bg-cream-50"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-farm-700 mb-1">
            Status
          </label>
          <select
            name="status"
            defaultValue="attending"
            className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 focus:border-transparent bg-cream-50"
          >
            <option value="attending">Attending</option>
            <option value="maybe">Maybe</option>
            <option value="declined">Declined</option>
          </select>
        </div>

        <div className="flex items-end gap-2">
          <button
            type="submit"
            disabled={loading}
            className="px-4 py-2 bg-forest-700 text-white rounded-md hover:bg-forest-800 disabled:opacity-50 transition-colors"
          >
            {loading ? 'Adding...' : 'Add'}
          </button>
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 border border-farm-400 text-farm-700 rounded-md hover:bg-farm-100 transition-colors"
          >
            Cancel
          </button>
        </div>
      </div>
    </form>
  )
}
