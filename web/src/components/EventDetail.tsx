import { useState } from 'react'
import type { EventWithMeals, CreateAttendeeRequest, Attendee, CreateMealRequest, CreateMealItemRequest } from '../api/types'
import { AttendeeForm } from './AttendeeForm'
import { MealSection } from './MealSection'

interface User {
  id: string
  email: string
  name: string
  picture: string
}

interface EventDetailProps {
  event: EventWithMeals
  currentUser: User | null
  onBack: () => void
  onEdit: () => void
  onAddAttendee: (data: CreateAttendeeRequest) => Promise<void>
  onUpdateAttendee: (attendeeId: string, status: string) => Promise<void>
  onRemoveAttendee: (attendeeId: string) => Promise<void>
  onCreateMeal: (data: CreateMealRequest) => Promise<void>
  onDeleteMeal: (mealId: string) => Promise<void>
  onAddMealItem: (mealId: string, data: CreateMealItemRequest) => Promise<void>
  onDeleteMealItem: (mealId: string, itemId: string) => Promise<void>
  onSignupForItem: (mealId: string, itemId: string, notes?: string) => Promise<void>
  onRemoveSignup: (mealId: string, itemId: string) => Promise<void>
}

function formatDateTime(dateString: string) {
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  })
}

const statusColors: Record<string, string> = {
  attending: 'bg-forest-100 text-forest-800',
  maybe: 'bg-farm-200 text-farm-800',
  declined: 'bg-rust-100 text-rust-800',
}

export function EventDetail({
  event,
  currentUser,
  onBack,
  onEdit,
  onAddAttendee,
  onUpdateAttendee,
  onRemoveAttendee,
  onCreateMeal,
  onDeleteMeal,
  onAddMealItem,
  onDeleteMealItem,
  onSignupForItem,
  onRemoveSignup,
}: EventDetailProps) {
  const [showAddForm, setShowAddForm] = useState(false)
  const [rsvpLoading, setRsvpLoading] = useState(false)
  const [linkCopied, setLinkCopied] = useState(false)

  const copyEventLink = async () => {
    const url = `${window.location.origin}/event/${event.id}`
    await navigator.clipboard.writeText(url)
    setLinkCopied(true)
    setTimeout(() => setLinkCopied(false), 2000)
  }

  const attendees = event.attendees || []
  const attending = attendees.filter((a) => a.status === 'attending')
  const maybe = attendees.filter((a) => a.status === 'maybe')
  const declined = attendees.filter((a) => a.status === 'declined')

  // Check if current user is already an attendee
  const myAttendee = currentUser
    ? attendees.find((a) => a.email.toLowerCase() === currentUser.email.toLowerCase())
    : null

  const handleAddAttendee = async (data: CreateAttendeeRequest) => {
    await onAddAttendee(data)
    setShowAddForm(false)
  }

  const handleRsvp = async (status: 'attending' | 'maybe' | 'declined') => {
    if (!currentUser) return
    setRsvpLoading(true)
    try {
      if (myAttendee) {
        await onUpdateAttendee(myAttendee.id, status)
      } else {
        await onAddAttendee({
          name: currentUser.name,
          email: currentUser.email,
          status,
        })
      }
    } finally {
      setRsvpLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      <button
        onClick={onBack}
        className="flex items-center gap-2 text-farm-600 hover:text-farm-900"
      >
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
        </svg>
        Back to events
      </button>

      <div className="bg-cream-50 rounded-lg shadow p-6 border border-farm-200">
        <div className="flex items-start justify-between gap-4">
          <h1 className="text-2xl font-bold text-farm-900">{event.title}</h1>
          <div className="flex items-center gap-2 shrink-0">
            <button
              onClick={onEdit}
              className="flex items-center gap-2 px-3 py-1.5 text-sm bg-forest-100 text-forest-700 rounded-md hover:bg-forest-200 transition-colors"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
              Edit
            </button>
            <button
              onClick={copyEventLink}
              className="flex items-center gap-2 px-3 py-1.5 text-sm bg-farm-100 text-farm-700 rounded-md hover:bg-farm-200 transition-colors"
            >
            {linkCopied ? (
              <>
                <svg className="w-4 h-4 text-forest-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                Copied!
              </>
            ) : (
              <>
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
                </svg>
                Share
              </>
            )}
            </button>
          </div>
        </div>

        {event.description && (
          <p className="mt-2 text-farm-600">{event.description}</p>
        )}

        <div className="mt-4 space-y-2 text-farm-600">
          <div className="flex items-center gap-2">
            <svg className="w-5 h-5 text-forest-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <span>{formatDateTime(event.start_time)}</span>
          </div>
          <div className="flex items-center gap-2">
            <svg className="w-5 h-5 text-forest-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>Until {formatDateTime(event.end_time)}</span>
          </div>
          {event.location && (
            <div className="flex items-center gap-2">
              <svg className="w-5 h-5 text-forest-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <span>{event.location}</span>
            </div>
          )}
        </div>
      </div>

      {/* RSVP Section */}
      {currentUser && (
        <div className="bg-cream-50 rounded-lg shadow p-6 border border-farm-200">
          <h2 className="text-lg font-semibold text-farm-800 mb-3">Your RSVP</h2>
          {myAttendee ? (
            <div className="flex items-center gap-4">
              <span className="text-farm-600">
                You're marked as <span className={`font-medium px-2 py-1 rounded ${statusColors[myAttendee.status]}`}>{myAttendee.status}</span>
              </span>
              <span className="text-farm-400">|</span>
              <span className="text-sm text-farm-500">Change to:</span>
              {myAttendee.status !== 'attending' && (
                <button
                  onClick={() => handleRsvp('attending')}
                  disabled={rsvpLoading}
                  className="px-3 py-1 text-sm bg-forest-600 text-white rounded hover:bg-forest-700 disabled:opacity-50 transition-colors"
                >
                  Attending
                </button>
              )}
              {myAttendee.status !== 'maybe' && (
                <button
                  onClick={() => handleRsvp('maybe')}
                  disabled={rsvpLoading}
                  className="px-3 py-1 text-sm bg-farm-500 text-white rounded hover:bg-farm-600 disabled:opacity-50 transition-colors"
                >
                  Maybe
                </button>
              )}
              {myAttendee.status !== 'declined' && (
                <button
                  onClick={() => handleRsvp('declined')}
                  disabled={rsvpLoading}
                  className="px-3 py-1 text-sm bg-rust-500 text-white rounded hover:bg-rust-600 disabled:opacity-50 transition-colors"
                >
                  Can't Go
                </button>
              )}
            </div>
          ) : (
            <div className="flex items-center gap-3">
              <span className="text-farm-600 mr-2">Will you attend?</span>
              <button
                onClick={() => handleRsvp('attending')}
                disabled={rsvpLoading}
                className="px-4 py-2 bg-forest-600 text-white rounded-md hover:bg-forest-700 disabled:opacity-50 transition-colors"
              >
                Yes, I'll be there
              </button>
              <button
                onClick={() => handleRsvp('maybe')}
                disabled={rsvpLoading}
                className="px-4 py-2 bg-farm-500 text-white rounded-md hover:bg-farm-600 disabled:opacity-50 transition-colors"
              >
                Maybe
              </button>
              <button
                onClick={() => handleRsvp('declined')}
                disabled={rsvpLoading}
                className="px-4 py-2 bg-rust-500 text-white rounded-md hover:bg-rust-600 disabled:opacity-50 transition-colors"
              >
                Can't Go
              </button>
            </div>
          )}
        </div>
      )}

      <div className="bg-cream-50 rounded-lg shadow p-6 border border-farm-200">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold text-farm-800">
            Attendees ({attendees.length})
          </h2>
          {!showAddForm && (
            <button
              onClick={() => setShowAddForm(true)}
              className="px-4 py-2 bg-rust-600 text-white rounded-md hover:bg-rust-700 transition-colors"
            >
              Add Attendee
            </button>
          )}
        </div>

        {showAddForm && (
          <div className="mb-6">
            <AttendeeForm
              onSubmit={handleAddAttendee}
              onCancel={() => setShowAddForm(false)}
            />
          </div>
        )}

        <div className="flex gap-4 mb-4 text-sm">
          <span className="px-2 py-1 bg-forest-100 text-forest-800 rounded">
            {attending.length} attending
          </span>
          <span className="px-2 py-1 bg-farm-200 text-farm-800 rounded">
            {maybe.length} maybe
          </span>
          <span className="px-2 py-1 bg-rust-100 text-rust-800 rounded">
            {declined.length} declined
          </span>
        </div>

        {attendees.length === 0 ? (
          <p className="text-farm-500 text-center py-4">
            No attendees yet. Add someone!
          </p>
        ) : (
          <div className="space-y-2">
            {attendees.map((attendee) => (
              <AttendeeRow
                key={attendee.id}
                attendee={attendee}
                onUpdateStatus={(status) => onUpdateAttendee(attendee.id, status)}
                onRemove={() => onRemoveAttendee(attendee.id)}
              />
            ))}
          </div>
        )}
      </div>

      {/* Meals Section */}
      <MealSection
        meals={event.meals || []}
        currentUser={currentUser}
        onCreateMeal={onCreateMeal}
        onDeleteMeal={onDeleteMeal}
        onAddItem={onAddMealItem}
        onDeleteItem={onDeleteMealItem}
        onSignup={onSignupForItem}
        onRemoveSignup={onRemoveSignup}
      />
    </div>
  )
}

interface AttendeeRowProps {
  attendee: Attendee
  onUpdateStatus: (status: string) => void
  onRemove: () => void
}

function AttendeeRow({ attendee, onUpdateStatus, onRemove }: AttendeeRowProps) {
  return (
    <div className="flex items-center justify-between p-3 bg-cream-200 rounded-lg border border-farm-200">
      <div className="flex items-center gap-3">
        <div className="w-10 h-10 bg-forest-200 rounded-full flex items-center justify-center text-forest-800 font-medium">
          {attendee.name.charAt(0).toUpperCase()}
        </div>
        <div>
          <div className="font-medium text-farm-900">{attendee.name}</div>
          <div className="text-sm text-farm-500">{attendee.email}</div>
        </div>
      </div>

      <div className="flex items-center gap-2">
        <select
          value={attendee.status}
          onChange={(e) => onUpdateStatus(e.target.value)}
          className={`px-3 py-1 rounded text-sm font-medium ${statusColors[attendee.status]}`}
        >
          <option value="attending">Attending</option>
          <option value="maybe">Maybe</option>
          <option value="declined">Declined</option>
        </select>

        <button
          onClick={onRemove}
          className="p-1 text-farm-400 hover:text-rust-600"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
  )
}
