import { useState } from 'react'
import type { CreateEventRequest, Event } from '../api/types'

interface EventFormProps {
  onSubmit: (data: CreateEventRequest) => Promise<void>
  onCancel: () => void
  event?: Event // Optional event for editing
}

// Helper to format date for input[type="date"] - uses local date
function formatDateForInput(dateString: string): string {
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// Helper to format time for input[type="time"] - uses local time
function formatTimeForInput(dateString: string): string {
  const date = new Date(dateString)
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${hours}:${minutes}`
}

// Format datetime with local timezone offset (e.g., 2025-12-20T09:00:00-05:00)
function formatLocalDateTime(date: string, time: string): string {
  const d = new Date(`${date}T${time}`)
  const tzOffset = -d.getTimezoneOffset()
  const sign = tzOffset >= 0 ? '+' : '-'
  const hours = String(Math.floor(Math.abs(tzOffset) / 60)).padStart(2, '0')
  const minutes = String(Math.abs(tzOffset) % 60).padStart(2, '0')
  return `${date}T${time}:00${sign}${hours}:${minutes}`
}

export function EventForm({ onSubmit, onCancel, event }: EventFormProps) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const isEditing = !!event

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    const form = e.currentTarget
    const formData = new FormData(form)

    const startDate = formData.get('start_date') as string
    const startTime = formData.get('start_time') as string
    const endDate = formData.get('end_date') as string
    const endTime = formData.get('end_time') as string

    const data: CreateEventRequest = {
      title: formData.get('title') as string,
      description: formData.get('description') as string,
      location: formData.get('location') as string,
      start_time: formatLocalDateTime(startDate, startTime),
      end_time: formatLocalDateTime(endDate, endTime),
    }

    try {
      await onSubmit(data)
      if (!isEditing) form.reset()
    } catch (err) {
      setError(err instanceof Error ? err.message : isEditing ? 'Failed to update event' : 'Failed to create event')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="bg-cream-50 rounded-lg shadow p-6 border border-farm-200">
      <h2 className="text-xl font-semibold mb-4 text-farm-800">
        {isEditing ? 'Edit Event' : 'Create New Event'}
      </h2>

      {error && (
        <div className="mb-4 p-3 bg-rust-100 text-rust-800 rounded border border-rust-300">{error}</div>
      )}

      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-farm-700 mb-1">
            Title *
          </label>
          <input
            type="text"
            name="title"
            required
            defaultValue={event?.title}
            className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 focus:border-transparent bg-cream-50"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-farm-700 mb-1">
            Description
          </label>
          <textarea
            name="description"
            rows={3}
            defaultValue={event?.description}
            className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 focus:border-transparent bg-cream-50"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-farm-700 mb-1">
            Location
          </label>
          <input
            type="text"
            name="location"
            defaultValue={event?.location}
            className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 focus:border-transparent bg-cream-50"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-farm-700 mb-1">
              Start Date *
            </label>
            <input
              type="date"
              name="start_date"
              required
              defaultValue={event ? formatDateForInput(event.start_time) : undefined}
              className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 focus:border-transparent bg-cream-50"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-farm-700 mb-1">
              Start Time *
            </label>
            <input
              type="time"
              name="start_time"
              required
              defaultValue={event ? formatTimeForInput(event.start_time) : undefined}
              className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 focus:border-transparent bg-cream-50"
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-farm-700 mb-1">
              End Date *
            </label>
            <input
              type="date"
              name="end_date"
              required
              defaultValue={event ? formatDateForInput(event.end_time) : undefined}
              className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 focus:border-transparent bg-cream-50"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-farm-700 mb-1">
              End Time *
            </label>
            <input
              type="time"
              name="end_time"
              required
              defaultValue={event ? formatTimeForInput(event.end_time) : undefined}
              className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 focus:border-transparent bg-cream-50"
            />
          </div>
        </div>
      </div>

      <div className="mt-6 flex gap-3">
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 bg-forest-700 text-white rounded-md hover:bg-forest-800 disabled:opacity-50 transition-colors"
        >
          {loading ? (isEditing ? 'Saving...' : 'Creating...') : (isEditing ? 'Save Changes' : 'Create Event')}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 border border-farm-400 text-farm-700 rounded-md hover:bg-farm-100 transition-colors"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
