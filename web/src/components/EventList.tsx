import type { Event } from '../api/types'

interface EventListProps {
  events: Event[]
  onSelect: (event: Event) => void
  onDelete: (id: string) => void
}

function formatDateTime(dateString: string) {
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  })
}

export function EventList({ events, onSelect, onDelete }: EventListProps) {
  if (events.length === 0) {
    return (
      <div className="bg-cream-50 rounded-lg shadow p-8 text-center text-farm-600 border border-farm-200">
        No events yet. Create your first event!
      </div>
    )
  }

  return (
    <div className="space-y-3">
      {events.map((event) => (
        <div
          key={event.id}
          className="bg-cream-50 rounded-lg shadow p-4 hover:shadow-md transition-shadow cursor-pointer border border-farm-200 hover:border-farm-300"
          onClick={() => onSelect(event)}
        >
          <div className="flex justify-between items-start">
            <div className="flex-1">
              <h3 className="font-semibold text-lg text-farm-900">
                {event.title}
              </h3>
              {event.description && (
                <p className="text-farm-600 text-sm mt-1 line-clamp-2">
                  {event.description}
                </p>
              )}
              <div className="mt-2 flex flex-wrap gap-3 text-sm text-farm-500">
                <span className="flex items-center gap-1">
                  <svg className="w-4 h-4 text-forest-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  {formatDateTime(event.start_time)}
                </span>
                {event.location && (
                  <span className="flex items-center gap-1">
                    <svg className="w-4 h-4 text-forest-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                    {event.location}
                  </span>
                )}
              </div>
            </div>
            <button
              onClick={(e) => {
                e.stopPropagation()
                onDelete(event.id)
              }}
              className="ml-4 p-2 text-farm-400 hover:text-rust-600 hover:bg-rust-50 rounded"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}
