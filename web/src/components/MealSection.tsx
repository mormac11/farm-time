import { useState } from 'react'
import type { MealWithItems, MealType, CreateMealRequest, CreateMealItemRequest } from '../api/types'

interface User {
  id: string
  email: string
  name: string
  picture: string
}

interface MealSectionProps {
  meals: MealWithItems[]
  currentUser: User | null
  onCreateMeal: (data: CreateMealRequest) => Promise<void>
  onDeleteMeal: (mealId: string) => Promise<void>
  onAddItem: (mealId: string, data: CreateMealItemRequest) => Promise<void>
  onDeleteItem: (mealId: string, itemId: string) => Promise<void>
  onSignup: (mealId: string, itemId: string, notes?: string) => Promise<void>
  onRemoveSignup: (mealId: string, itemId: string) => Promise<void>
}

const mealTypeLabels: Record<MealType, string> = {
  breakfast: 'Breakfast',
  lunch: 'Lunch',
  dinner: 'Dinner',
  snacks: 'Snacks',
  other: 'Other',
}

const mealTypeIcons: Record<MealType, string> = {
  breakfast: '🌅',
  lunch: '☀️',
  dinner: '🍽️',
  snacks: '🍿',
  other: '🍴',
}

export function MealSection({
  meals,
  currentUser,
  onCreateMeal,
  onDeleteMeal,
  onAddItem,
  onDeleteItem,
  onSignup,
  onRemoveSignup,
}: MealSectionProps) {
  const [showMealForm, setShowMealForm] = useState(false)
  const [addingItemToMeal, setAddingItemToMeal] = useState<string | null>(null)
  const [signupItemId, setSignupItemId] = useState<string | null>(null)
  const [signupMealId, setSignupMealId] = useState<string | null>(null)
  const [signupNotes, setSignupNotes] = useState('')
  const [loading, setLoading] = useState(false)

  const handleCreateMeal = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setLoading(true)
    const form = e.currentTarget
    const formData = new FormData(form)

    const mealType = formData.get('meal_type') as MealType
    let name = formData.get('name') as string
    if (!name) {
      name = mealTypeLabels[mealType]
    }

    try {
      await onCreateMeal({
        name,
        meal_type: mealType,
        meal_date: formData.get('meal_date') as string || undefined,
        notes: formData.get('notes') as string || undefined,
      })
      form.reset()
      setShowMealForm(false)
    } finally {
      setLoading(false)
    }
  }

  const handleAddItem = async (e: React.FormEvent<HTMLFormElement>, mealId: string) => {
    e.preventDefault()
    setLoading(true)
    const form = e.currentTarget
    const formData = new FormData(form)

    try {
      await onAddItem(mealId, {
        name: formData.get('name') as string,
        description: formData.get('description') as string || undefined,
      })
      form.reset()
      setAddingItemToMeal(null)
    } finally {
      setLoading(false)
    }
  }

  const handleSignup = async () => {
    if (!signupItemId || !signupMealId) return
    setLoading(true)
    try {
      await onSignup(signupMealId, signupItemId, signupNotes || undefined)
      setSignupItemId(null)
      setSignupMealId(null)
      setSignupNotes('')
    } finally {
      setLoading(false)
    }
  }

  const isUserSignedUp = (itemId: string) => {
    if (!currentUser) return false
    for (const meal of meals) {
      for (const item of meal.items) {
        if (item.id === itemId) {
          return item.signups.some(s => s.user_id === currentUser.id)
        }
      }
    }
    return false
  }

  return (
    <div className="bg-cream-50 rounded-lg shadow p-6 border border-farm-200">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold text-farm-800">
          Meals ({meals.length})
        </h2>
        {!showMealForm && (
          <button
            onClick={() => setShowMealForm(true)}
            className="px-4 py-2 bg-rust-600 text-white rounded-md hover:bg-rust-700 transition-colors"
          >
            + Add Meal
          </button>
        )}
      </div>

      {/* Add Meal Form */}
      {showMealForm && (
        <form onSubmit={handleCreateMeal} className="mb-6 p-4 bg-cream-200 rounded-lg border border-farm-200">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-farm-700 mb-1">
                Meal Type *
              </label>
              <select
                name="meal_type"
                required
                className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 bg-cream-50"
              >
                <option value="breakfast">Breakfast</option>
                <option value="lunch">Lunch</option>
                <option value="dinner">Dinner</option>
                <option value="snacks">Snacks</option>
                <option value="other">Other</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-farm-700 mb-1">
                Custom Name (optional)
              </label>
              <input
                type="text"
                name="name"
                placeholder="e.g., Saturday Dinner"
                className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 bg-cream-50"
              />
            </div>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-farm-700 mb-1">
                Date (optional)
              </label>
              <input
                type="date"
                name="meal_date"
                className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 bg-cream-50"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-farm-700 mb-1">
                Notes (optional)
              </label>
              <input
                type="text"
                name="notes"
                placeholder="e.g., Vegetarian options available"
                className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 bg-cream-50"
              />
            </div>
          </div>
          <div className="flex gap-2">
            <button
              type="submit"
              disabled={loading}
              className="px-4 py-2 bg-forest-700 text-white rounded-md hover:bg-forest-800 disabled:opacity-50 transition-colors"
            >
              {loading ? 'Adding...' : 'Add Meal'}
            </button>
            <button
              type="button"
              onClick={() => setShowMealForm(false)}
              className="px-4 py-2 border border-farm-400 text-farm-700 rounded-md hover:bg-farm-100 transition-colors"
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      {/* Meals List */}
      {meals.length === 0 ? (
        <p className="text-farm-500 text-center py-4">
          No meals planned yet. Add a meal to start organizing!
        </p>
      ) : (
        <div className="space-y-4">
          {meals.map((meal) => (
            <div key={meal.id} className="border border-farm-200 rounded-lg p-4 bg-cream-100">
              <div className="flex justify-between items-start mb-3">
                <div className="flex items-center gap-2">
                  <span className="text-xl">{mealTypeIcons[meal.meal_type as MealType]}</span>
                  <div>
                    <h3 className="font-semibold text-farm-900">{meal.name}</h3>
                    {meal.meal_date && (
                      <span className="text-sm text-farm-500">
                        {new Date(meal.meal_date + 'T00:00:00').toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' })}
                      </span>
                    )}
                  </div>
                </div>
                <button
                  onClick={() => onDeleteMeal(meal.id)}
                  className="p-1 text-farm-400 hover:text-rust-600"
                  title="Delete meal"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              {meal.notes && (
                <p className="text-sm text-farm-600 mb-3">{meal.notes}</p>
              )}

              {/* Items */}
              <div className="space-y-2">
                {meal.items.map((item) => {
                  const userSignedUp = isUserSignedUp(item.id)
                  return (
                    <div key={item.id} className="flex items-start justify-between p-2 bg-cream-50 rounded border border-farm-100">
                      <div className="flex-1">
                        <div className="flex items-center gap-2">
                          <span className="font-medium text-farm-800">{item.name}</span>
                          <button
                            onClick={() => onDeleteItem(meal.id, item.id)}
                            className="p-0.5 text-farm-300 hover:text-rust-500"
                            title="Delete item"
                          >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                          </button>
                        </div>
                        {item.description && (
                          <p className="text-sm text-farm-500">{item.description}</p>
                        )}
                        {item.signups.length > 0 ? (
                          <div className="mt-1 space-y-1">
                            {item.signups.map((signup) => (
                              <div key={signup.id} className="flex items-center gap-2 text-sm">
                                <span className="text-forest-600">
                                  {signup.user_id === currentUser?.id ? 'You' : signup.user_name}
                                </span>
                                {signup.notes && (
                                  <span className="text-farm-500">- "{signup.notes}"</span>
                                )}
                                {signup.user_id === currentUser?.id && (
                                  <button
                                    onClick={() => onRemoveSignup(meal.id, item.id)}
                                    className="text-rust-500 hover:text-rust-700 text-xs"
                                  >
                                    (remove)
                                  </button>
                                )}
                              </div>
                            ))}
                          </div>
                        ) : (
                          <p className="text-sm text-farm-400 mt-1">No one signed up yet</p>
                        )}
                      </div>
                      {currentUser && !userSignedUp && (
                        <button
                          onClick={() => {
                            setSignupItemId(item.id)
                            setSignupMealId(meal.id)
                          }}
                          className="px-3 py-1 text-sm bg-forest-100 text-forest-700 rounded hover:bg-forest-200 transition-colors whitespace-nowrap"
                        >
                          I'll bring this
                        </button>
                      )}
                      {currentUser && userSignedUp && (
                        <span className="px-3 py-1 text-sm bg-forest-600 text-white rounded">
                          Signed up
                        </span>
                      )}
                    </div>
                  )
                })}
              </div>

              {/* Add Item Form */}
              {addingItemToMeal === meal.id ? (
                <form onSubmit={(e) => handleAddItem(e, meal.id)} className="mt-3 p-3 bg-cream-200 rounded border border-farm-200">
                  <div className="flex gap-2">
                    <input
                      type="text"
                      name="name"
                      required
                      placeholder="Item name (e.g., Burgers)"
                      className="flex-1 px-3 py-1.5 text-sm border border-farm-300 rounded focus:outline-none focus:ring-2 focus:ring-forest-500 bg-cream-50"
                    />
                    <input
                      type="text"
                      name="description"
                      placeholder="Description (optional)"
                      className="flex-1 px-3 py-1.5 text-sm border border-farm-300 rounded focus:outline-none focus:ring-2 focus:ring-forest-500 bg-cream-50"
                    />
                    <button
                      type="submit"
                      disabled={loading}
                      className="px-3 py-1.5 text-sm bg-forest-700 text-white rounded hover:bg-forest-800 disabled:opacity-50"
                    >
                      Add
                    </button>
                    <button
                      type="button"
                      onClick={() => setAddingItemToMeal(null)}
                      className="px-3 py-1.5 text-sm border border-farm-400 text-farm-700 rounded hover:bg-farm-100"
                    >
                      Cancel
                    </button>
                  </div>
                </form>
              ) : (
                <button
                  onClick={() => setAddingItemToMeal(meal.id)}
                  className="mt-3 text-sm text-forest-600 hover:text-forest-800 flex items-center gap-1"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                  </svg>
                  Add Item
                </button>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Signup Modal */}
      {signupItemId && signupMealId && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-cream-50 rounded-lg p-6 max-w-md w-full mx-4 shadow-xl">
            <h3 className="text-lg font-semibold text-farm-800 mb-4">Sign up to bring this item</h3>
            <div className="mb-4">
              <label className="block text-sm font-medium text-farm-700 mb-1">
                Notes (optional)
              </label>
              <input
                type="text"
                value={signupNotes}
                onChange={(e) => setSignupNotes(e.target.value)}
                placeholder="e.g., I'll bring 2 pounds"
                className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 bg-cream-50"
              />
            </div>
            <div className="flex gap-3">
              <button
                onClick={handleSignup}
                disabled={loading}
                className="flex-1 px-4 py-2 bg-forest-700 text-white rounded-md hover:bg-forest-800 disabled:opacity-50 transition-colors"
              >
                {loading ? 'Signing up...' : 'Sign Up'}
              </button>
              <button
                onClick={() => {
                  setSignupItemId(null)
                  setSignupMealId(null)
                  setSignupNotes('')
                }}
                className="px-4 py-2 border border-farm-400 text-farm-700 rounded-md hover:bg-farm-100 transition-colors"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
