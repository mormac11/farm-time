import { useState, useEffect, useCallback } from 'react'
import { api } from './api/client'
import type { Event, EventWithAll, CreateEventRequest, CreateAttendeeRequest, CreateMealRequest, CreateMealItemRequest, CreateTodoRequest, UpdateTodoRequest } from './api/types'
import { EventList } from './components/EventList'
import { EventForm } from './components/EventForm'
import { EventDetail } from './components/EventDetail'
import { useAuth } from './contexts/AuthContext'

type View = 'list' | 'create' | 'detail' | 'edit'

// Get event ID from URL path like /event/abc123
function getEventIdFromUrl(): string | null {
  const match = window.location.pathname.match(/^\/event\/([^/]+)/)
  return match ? match[1] : null
}

function App() {
  const { user, loading: authLoading, login, logout } = useAuth()
  const [view, setView] = useState<View>('list')
  const [events, setEvents] = useState<Event[]>([])
  const [selectedEvent, setSelectedEvent] = useState<EventWithAll | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const loadEvents = useCallback(async () => {
    try {
      const data = await api.listEvents()
      setEvents(data || [])
      setError('')
    } catch (err) {
      setError('Failed to load events')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  // Load event from URL on mount
  useEffect(() => {
    const eventId = getEventIdFromUrl()
    if (eventId && user) {
      api.getEvent(eventId)
        .then((fullEvent) => {
          setSelectedEvent(fullEvent)
          setView('detail')
        })
        .catch(() => {
          // Event not found, go to list
          window.history.replaceState(null, '', '/')
        })
    }
  }, [user])

  useEffect(() => {
    loadEvents()
  }, [loadEvents])

  const handleCreateEvent = async (data: CreateEventRequest) => {
    await api.createEvent(data)
    await loadEvents()
    setView('list')
  }

  const handleUpdateEvent = async (data: CreateEventRequest) => {
    if (!selectedEvent) return
    await api.updateEvent(selectedEvent.id, data)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
    setView('detail')
  }

  const handleEditEvent = () => {
    setView('edit')
  }

  const handleDeleteEvent = async (id: string) => {
    if (!confirm('Are you sure you want to delete this event?')) return
    await api.deleteEvent(id)
    await loadEvents()
  }

  const handleSelectEvent = async (event: Event) => {
    const fullEvent = await api.getEvent(event.id)
    setSelectedEvent(fullEvent)
    setView('detail')
    window.history.pushState(null, '', `/event/${event.id}`)
  }

  const handleAddAttendee = async (data: CreateAttendeeRequest) => {
    if (!selectedEvent) return
    await api.addAttendee(selectedEvent.id, data)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  const handleUpdateAttendee = async (attendeeId: string, status: string) => {
    if (!selectedEvent) return
    await api.updateAttendee(selectedEvent.id, attendeeId, { status: status as CreateAttendeeRequest['status'] })
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  const handleRemoveAttendee = async (attendeeId: string) => {
    if (!selectedEvent) return
    await api.removeAttendee(selectedEvent.id, attendeeId)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  // Meal handlers
  const handleCreateMeal = async (data: CreateMealRequest) => {
    if (!selectedEvent) return
    await api.createMeal(selectedEvent.id, data)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  const handleDeleteMeal = async (mealId: string) => {
    if (!selectedEvent) return
    await api.deleteMeal(selectedEvent.id, mealId)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  const handleAddMealItem = async (mealId: string, data: CreateMealItemRequest) => {
    if (!selectedEvent) return
    await api.addMealItem(selectedEvent.id, mealId, data)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  const handleDeleteMealItem = async (mealId: string, itemId: string) => {
    if (!selectedEvent) return
    await api.deleteMealItem(selectedEvent.id, mealId, itemId)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  const handleSignupForItem = async (mealId: string, itemId: string, notes?: string) => {
    if (!selectedEvent) return
    await api.signupForItem(selectedEvent.id, mealId, itemId, notes ? { notes } : undefined)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  const handleRemoveSignup = async (mealId: string, itemId: string) => {
    if (!selectedEvent) return
    await api.removeSignup(selectedEvent.id, mealId, itemId)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  // Todo handlers
  const handleCreateTodo = async (data: CreateTodoRequest) => {
    if (!selectedEvent) return
    await api.createTodo(selectedEvent.id, data)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  const handleUpdateTodo = async (todoId: string, data: UpdateTodoRequest) => {
    if (!selectedEvent) return
    await api.updateTodo(selectedEvent.id, todoId, data)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  const handleDeleteTodo = async (todoId: string) => {
    if (!selectedEvent) return
    await api.deleteTodo(selectedEvent.id, todoId)
    const updated = await api.getEvent(selectedEvent.id)
    setSelectedEvent(updated)
  }

  const handleBack = () => {
    setSelectedEvent(null)
    setView('list')
    loadEvents()
    window.history.pushState(null, '', '/')
  }

  // Tractor icon for rustic farm theme
  const TractorIcon = ({ className }: { className?: string }) => (
    <svg className={className} viewBox="0 0 24 24" fill="currentColor">
      <path d="M19 13a3 3 0 1 0 0 6 3 3 0 0 0 0-6zm0 4.5a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3z"/>
      <path d="M7 14a4 4 0 1 0 0 8 4 4 0 0 0 0-8zm0 6a2 2 0 1 1 0-4 2 2 0 0 1 0 4z"/>
      <path d="M20 10h-3V7h-2V5H9v2H7v1H5v2H3v3h1.05A5 5 0 0 1 7 11c1.53 0 2.92.58 4 1.53V10h5v1.29A4 4 0 0 1 19 10h1zm-9 0H9V7h2v3z"/>
    </svg>
  )

  // Show loading spinner while checking auth
  if (authLoading) {
    return (
      <div className="min-h-screen bg-cream-100 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-4 border-forest-600 border-t-transparent"></div>
      </div>
    )
  }

  // Show login page if not authenticated
  if (!user) {
    return (
      <div className="min-h-screen bg-cream-100 flex flex-col items-center justify-center">
        <div className="text-center">
          <div className="flex items-center justify-center gap-3 mb-4">
            <TractorIcon className="w-14 h-14 text-forest-700" />
            <h1 className="text-4xl font-bold text-farm-800">Farm Time</h1>
          </div>
          <p className="text-farm-600 mb-8 text-lg">Event Scheduler & Attendance Tracker</p>
          <button
            onClick={login}
            className="px-6 py-3 bg-forest-700 text-white rounded-lg hover:bg-forest-800 flex items-center gap-3 mx-auto transition-colors shadow-md"
          >
            <svg className="w-5 h-5" viewBox="0 0 24 24">
              <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="currentColor" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="currentColor" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="currentColor" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            Sign in with Google
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-cream-100">
      <header className="bg-forest-800 text-white shadow-lg">
        <div className="max-w-4xl mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <button
              onClick={handleBack}
              className="text-left hover:opacity-80 transition-opacity"
            >
              <div className="flex items-center gap-3">
                <TractorIcon className="w-8 h-8 text-cream-200" />
                <h1 className="text-2xl font-bold">Farm Time</h1>
              </div>
              <p className="mt-1 text-forest-200">Event Scheduler & Attendance Tracker</p>
            </button>
            <div className="flex items-center gap-4">
              <img
                src={user.picture}
                alt={user.name}
                className="w-8 h-8 rounded-full border-2 border-cream-200"
                referrerPolicy="no-referrer"
              />
              <span className="hidden sm:inline text-cream-100">{user.name}</span>
              <button
                onClick={logout}
                className="text-sm text-forest-200 hover:text-white underline"
              >
                Sign out
              </button>
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 py-8">
        {error && (
          <div className="mb-6 p-4 bg-rust-100 text-rust-800 rounded-lg border border-rust-300">
            {error}
          </div>
        )}

        {loading ? (
          <div className="flex justify-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-4 border-forest-600 border-t-transparent"></div>
          </div>
        ) : (
          <>
            {view === 'list' && (
              <div>
                <div className="flex justify-between items-center mb-6">
                  <h2 className="text-xl font-semibold text-farm-800">
                    Upcoming Events
                  </h2>
                  <button
                    onClick={() => setView('create')}
                    className="px-4 py-2 bg-rust-600 text-white rounded-md hover:bg-rust-700 flex items-center gap-2 shadow-sm transition-colors"
                  >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                    </svg>
                    New Event
                  </button>
                </div>
                <EventList
                  events={events}
                  onSelect={handleSelectEvent}
                  onDelete={handleDeleteEvent}
                />
              </div>
            )}

            {view === 'create' && (
              <EventForm
                onSubmit={handleCreateEvent}
                onCancel={() => setView('list')}
              />
            )}

            {view === 'detail' && selectedEvent && (
              <EventDetail
                event={selectedEvent}
                currentUser={user}
                onBack={handleBack}
                onEdit={handleEditEvent}
                onAddAttendee={handleAddAttendee}
                onUpdateAttendee={handleUpdateAttendee}
                onRemoveAttendee={handleRemoveAttendee}
                onCreateMeal={handleCreateMeal}
                onDeleteMeal={handleDeleteMeal}
                onAddMealItem={handleAddMealItem}
                onDeleteMealItem={handleDeleteMealItem}
                onSignupForItem={handleSignupForItem}
                onRemoveSignup={handleRemoveSignup}
                onCreateTodo={handleCreateTodo}
                onUpdateTodo={handleUpdateTodo}
                onDeleteTodo={handleDeleteTodo}
              />
            )}

            {view === 'edit' && selectedEvent && (
              <EventForm
                onSubmit={handleUpdateEvent}
                onCancel={() => setView('detail')}
                event={selectedEvent}
              />
            )}
          </>
        )}
      </main>

      <footer className="mt-auto py-6 text-center text-farm-500 text-sm">
        Farm Time - Event Scheduler
      </footer>
    </div>
  )
}

export default App
