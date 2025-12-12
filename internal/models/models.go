package models

import "time"

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Attendee struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"` // "attending", "maybe", "declined"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EventWithAttendees struct {
	Event
	Attendees []Attendee `json:"attendees"`
}

type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type UpdateEventRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Location    *string    `json:"location,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
}

type CreateAttendeeRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type UpdateAttendeeRequest struct {
	Name   *string `json:"name,omitempty"`
	Email  *string `json:"email,omitempty"`
	Status *string `json:"status,omitempty"`
}

// User represents a user authenticated via Google OAuth
type User struct {
	ID              string    `json:"id"`
	GoogleID        string    `json:"-"` // Never expose to frontend
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	Picture         string    `json:"picture"`
	IsAdmin         bool      `json:"is_admin"`
	CanCreateEvents bool      `json:"can_create_events"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UpdateUserPermissionsRequest struct {
	CanCreateEvents *bool `json:"can_create_events,omitempty"`
}

// Session represents a user session
type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// Meal represents a meal within an event (e.g., "Saturday Dinner")
type Meal struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	Name      string    `json:"name"`
	MealType  string    `json:"meal_type"` // "breakfast", "lunch", "dinner", "snacks", "other"
	MealDate  *string   `json:"meal_date"` // Optional YYYY-MM-DD
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MealItem represents something needed for a meal (e.g., "Burgers")
type MealItem struct {
	ID                   string    `json:"id"`
	MealID               string    `json:"meal_id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	AssignedAttendeeID   *string   `json:"assigned_attendee_id"`
	AssignedAttendeeName *string   `json:"assigned_attendee_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// MealSignup represents a user signing up to bring a meal item
type MealSignup struct {
	ID         string    `json:"id"`
	MealItemID string    `json:"meal_item_id"`
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name"`
	UserEmail  string    `json:"user_email"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
}

// MealItemWithSignups includes item details with who's bringing it
type MealItemWithSignups struct {
	MealItem
	Signups []MealSignup `json:"signups"`
}

// MealWithItems includes meal details with all items
type MealWithItems struct {
	Meal
	Items []MealItemWithSignups `json:"items"`
}

// EventWithMeals extends EventWithAttendees to include meals
type EventWithMeals struct {
	EventWithAttendees
	Meals []MealWithItems `json:"meals"`
}

// Request types for meals
type CreateMealRequest struct {
	Name     string  `json:"name"`
	MealType string  `json:"meal_type"`
	MealDate *string `json:"meal_date"`
	Notes    string  `json:"notes"`
}

type UpdateMealRequest struct {
	Name     *string `json:"name,omitempty"`
	MealType *string `json:"meal_type,omitempty"`
	MealDate *string `json:"meal_date,omitempty"`
	Notes    *string `json:"notes,omitempty"`
}

type CreateMealItemRequest struct {
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	AssignedAttendeeID *string `json:"assigned_attendee_id"`
}

type UpdateMealItemRequest struct {
	Name               *string `json:"name,omitempty"`
	Description        *string `json:"description,omitempty"`
	AssignedAttendeeID *string `json:"assigned_attendee_id,omitempty"`
}

type CreateMealSignupRequest struct {
	Notes string `json:"notes"`
}

// Todo represents a task item for an event
type Todo struct {
	ID                   string    `json:"id"`
	EventID              string    `json:"event_id"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	Completed            bool      `json:"completed"`
	AssignedAttendeeID   *string   `json:"assigned_attendee_id"`
	AssignedAttendeeName *string   `json:"assigned_attendee_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type CreateTodoRequest struct {
	Title              string  `json:"title"`
	Description        string  `json:"description"`
	AssignedAttendeeID *string `json:"assigned_attendee_id"`
}

type UpdateTodoRequest struct {
	Title              *string `json:"title,omitempty"`
	Description        *string `json:"description,omitempty"`
	Completed          *bool   `json:"completed,omitempty"`
	AssignedAttendeeID *string `json:"assigned_attendee_id,omitempty"`
}

// EventWithAll extends EventWithMeals to include todos
type EventWithAll struct {
	EventWithMeals
	Todos []Todo `json:"todos"`
}
