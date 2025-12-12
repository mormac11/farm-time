import type {
  Event,
  EventWithAll,
  Attendee,
  CreateEventRequest,
  CreateAttendeeRequest,
  Meal,
  MealItem,
  MealSignup,
  MealWithItems,
  CreateMealRequest,
  CreateMealItemRequest,
  CreateMealSignupRequest,
  Todo,
  CreateTodoRequest,
  UpdateTodoRequest,
} from './types'

const API_BASE = '/api'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    ...options,
  })

  if (response.status === 401) {
    window.location.href = '/auth/google/login'
    throw new Error('Unauthorized')
  }

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }))
    throw new Error(error.error || 'Request failed')
  }

  if (response.status === 204) {
    return undefined as T
  }

  return response.json()
}

export const api = {
  // Events
  listEvents: () => request<Event[]>('/events'),

  getEvent: (id: string) => request<EventWithAll>(`/events/${id}`),

  createEvent: (data: CreateEventRequest) =>
    request<Event>('/events', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  updateEvent: (id: string, data: Partial<CreateEventRequest>) =>
    request<Event>(`/events/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  deleteEvent: (id: string) =>
    request<void>(`/events/${id}`, {
      method: 'DELETE',
    }),

  // Attendees
  listAttendees: (eventId: string) =>
    request<Attendee[]>(`/events/${eventId}/attendees`),

  addAttendee: (eventId: string, data: CreateAttendeeRequest) =>
    request<Attendee>(`/events/${eventId}/attendees`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  updateAttendee: (eventId: string, attendeeId: string, data: Partial<CreateAttendeeRequest>) =>
    request<Attendee>(`/events/${eventId}/attendees/${attendeeId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  removeAttendee: (eventId: string, attendeeId: string) =>
    request<void>(`/events/${eventId}/attendees/${attendeeId}`, {
      method: 'DELETE',
    }),

  // Meals
  listMeals: (eventId: string) =>
    request<MealWithItems[]>(`/events/${eventId}/meals`),

  createMeal: (eventId: string, data: CreateMealRequest) =>
    request<Meal>(`/events/${eventId}/meals`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  updateMeal: (eventId: string, mealId: string, data: Partial<CreateMealRequest>) =>
    request<Meal>(`/events/${eventId}/meals/${mealId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  deleteMeal: (eventId: string, mealId: string) =>
    request<void>(`/events/${eventId}/meals/${mealId}`, {
      method: 'DELETE',
    }),

  // Meal Items
  addMealItem: (eventId: string, mealId: string, data: CreateMealItemRequest) =>
    request<MealItem>(`/events/${eventId}/meals/${mealId}/items`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  updateMealItem: (eventId: string, mealId: string, itemId: string, data: Partial<CreateMealItemRequest>) =>
    request<MealItem>(`/events/${eventId}/meals/${mealId}/items/${itemId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  deleteMealItem: (eventId: string, mealId: string, itemId: string) =>
    request<void>(`/events/${eventId}/meals/${mealId}/items/${itemId}`, {
      method: 'DELETE',
    }),

  // Meal Signups
  signupForItem: (eventId: string, mealId: string, itemId: string, data?: CreateMealSignupRequest) =>
    request<MealSignup>(`/events/${eventId}/meals/${mealId}/items/${itemId}/signup`, {
      method: 'POST',
      body: JSON.stringify(data || {}),
    }),

  removeSignup: (eventId: string, mealId: string, itemId: string) =>
    request<void>(`/events/${eventId}/meals/${mealId}/items/${itemId}/signup`, {
      method: 'DELETE',
    }),

  // Todos
  listTodos: (eventId: string) =>
    request<Todo[]>(`/events/${eventId}/todos`),

  createTodo: (eventId: string, data: CreateTodoRequest) =>
    request<Todo>(`/events/${eventId}/todos`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  updateTodo: (eventId: string, todoId: string, data: UpdateTodoRequest) =>
    request<Todo>(`/events/${eventId}/todos/${todoId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  deleteTodo: (eventId: string, todoId: string) =>
    request<void>(`/events/${eventId}/todos/${todoId}`, {
      method: 'DELETE',
    }),
}
