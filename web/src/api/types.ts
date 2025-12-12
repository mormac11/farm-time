export interface Event {
  id: string
  title: string
  description: string
  location: string
  start_time: string
  end_time: string
  created_at: string
  updated_at: string
}

export interface Attendee {
  id: string
  event_id: string
  name: string
  email: string
  status: 'attending' | 'maybe' | 'declined'
  created_at: string
  updated_at: string
}

export interface EventWithAttendees extends Event {
  attendees: Attendee[]
}

export interface CreateEventRequest {
  title: string
  description: string
  location: string
  start_time: string
  end_time: string
}

export interface CreateAttendeeRequest {
  name: string
  email: string
  status: 'attending' | 'maybe' | 'declined'
}

// Meal types
export type MealType = 'breakfast' | 'lunch' | 'dinner' | 'snacks' | 'other'

export interface Meal {
  id: string
  event_id: string
  name: string
  meal_type: MealType
  meal_date: string | null
  notes: string
  created_at: string
  updated_at: string
}

export interface MealItem {
  id: string
  meal_id: string
  name: string
  description: string
  assigned_attendee_id: string | null
  assigned_attendee_name: string | null
  created_at: string
  updated_at: string
}

export interface MealSignup {
  id: string
  meal_item_id: string
  user_id: string
  user_name: string
  user_email: string
  notes: string
  created_at: string
}

export interface MealItemWithSignups extends MealItem {
  signups: MealSignup[]
}

export interface MealWithItems extends Meal {
  items: MealItemWithSignups[]
}

export interface EventWithMeals extends EventWithAttendees {
  meals: MealWithItems[]
}

export interface CreateMealRequest {
  name: string
  meal_type: MealType
  meal_date?: string
  notes?: string
}

export interface CreateMealItemRequest {
  name: string
  description?: string
  assigned_attendee_id?: string
}

export interface CreateMealSignupRequest {
  notes?: string
}
