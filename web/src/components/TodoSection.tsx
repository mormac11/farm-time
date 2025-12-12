import { useState } from 'react'
import type { Todo, CreateTodoRequest, Attendee } from '../api/types'

interface TodoSectionProps {
  todos: Todo[]
  attendees: Attendee[]
  onCreateTodo: (data: CreateTodoRequest) => Promise<void>
  onUpdateTodo: (todoId: string, data: { completed?: boolean; title?: string; description?: string; assigned_attendee_id?: string }) => Promise<void>
  onDeleteTodo: (todoId: string) => Promise<void>
}

export function TodoSection({
  todos,
  attendees,
  onCreateTodo,
  onUpdateTodo,
  onDeleteTodo,
}: TodoSectionProps) {
  const [showForm, setShowForm] = useState(false)
  const [loading, setLoading] = useState(false)

  const handleCreate = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setLoading(true)
    const form = e.currentTarget
    const formData = new FormData(form)

    const assignedAttendeeId = formData.get('assigned_attendee_id') as string

    try {
      await onCreateTodo({
        title: formData.get('title') as string,
        description: formData.get('description') as string || undefined,
        assigned_attendee_id: assignedAttendeeId || undefined,
      })
      form.reset()
      setShowForm(false)
    } finally {
      setLoading(false)
    }
  }

  const handleToggle = async (todo: Todo) => {
    await onUpdateTodo(todo.id, { completed: !todo.completed })
  }

  const completedTodos = todos.filter(t => t.completed)
  const pendingTodos = todos.filter(t => !t.completed)

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold text-farm-800">
          Tasks ({pendingTodos.length} pending, {completedTodos.length} done)
        </h3>
        {!showForm && (
          <button
            onClick={() => setShowForm(true)}
            className="px-4 py-2 bg-rust-600 text-white rounded-md hover:bg-rust-700 transition-colors"
          >
            + Add Task
          </button>
        )}
      </div>

      {/* Add Task Form */}
      {showForm && (
        <form onSubmit={handleCreate} className="p-4 bg-cream-200 rounded-lg border border-farm-200">
          <div className="space-y-3">
            <div>
              <label className="block text-sm font-medium text-farm-700 mb-1">
                Task *
              </label>
              <input
                type="text"
                name="title"
                required
                placeholder="What needs to be done?"
                className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 bg-cream-50"
              />
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div>
                <label className="block text-sm font-medium text-farm-700 mb-1">
                  Description (optional)
                </label>
                <input
                  type="text"
                  name="description"
                  placeholder="Additional details..."
                  className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 bg-cream-50"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-farm-700 mb-1">
                  Assign to (optional)
                </label>
                <select
                  name="assigned_attendee_id"
                  className="w-full px-3 py-2 border border-farm-300 rounded-md focus:outline-none focus:ring-2 focus:ring-forest-500 bg-cream-50"
                >
                  <option value="">Unassigned</option>
                  {attendees.map((attendee) => (
                    <option key={attendee.id} value={attendee.id}>
                      {attendee.name}
                    </option>
                  ))}
                </select>
              </div>
            </div>
          </div>
          <div className="flex gap-2 mt-4">
            <button
              type="submit"
              disabled={loading}
              className="px-4 py-2 bg-forest-700 text-white rounded-md hover:bg-forest-800 disabled:opacity-50 transition-colors"
            >
              {loading ? 'Adding...' : 'Add Task'}
            </button>
            <button
              type="button"
              onClick={() => setShowForm(false)}
              className="px-4 py-2 border border-farm-400 text-farm-700 rounded-md hover:bg-farm-100 transition-colors"
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      {/* Tasks List */}
      {todos.length === 0 ? (
        <p className="text-farm-500 text-center py-8">
          No tasks yet. Add a task to get started!
        </p>
      ) : (
        <div className="space-y-2">
          {/* Pending Tasks */}
          {pendingTodos.map((todo) => (
            <div
              key={todo.id}
              className="flex items-start gap-3 p-3 bg-cream-100 rounded-lg border border-farm-200"
            >
              <button
                onClick={() => handleToggle(todo)}
                className="mt-0.5 w-5 h-5 rounded border-2 border-farm-400 hover:border-forest-600 flex-shrink-0 flex items-center justify-center"
              >
                {/* Empty checkbox */}
              </button>
              <div className="flex-1 min-w-0">
                <p className="font-medium text-farm-800">{todo.title}</p>
                {todo.description && (
                  <p className="text-sm text-farm-500 mt-0.5">{todo.description}</p>
                )}
                {todo.assigned_attendee_name && (
                  <p className="text-sm text-forest-600 mt-1">
                    Assigned to: {todo.assigned_attendee_name}
                  </p>
                )}
              </div>
              <button
                onClick={() => onDeleteTodo(todo.id)}
                className="p-1 text-farm-400 hover:text-rust-600 flex-shrink-0"
                title="Delete task"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          ))}

          {/* Completed Tasks */}
          {completedTodos.length > 0 && (
            <div className="mt-4">
              <p className="text-sm text-farm-500 mb-2">Completed</p>
              {completedTodos.map((todo) => (
                <div
                  key={todo.id}
                  className="flex items-start gap-3 p-3 bg-cream-50 rounded-lg border border-farm-100 opacity-60"
                >
                  <button
                    onClick={() => handleToggle(todo)}
                    className="mt-0.5 w-5 h-5 rounded border-2 border-forest-600 bg-forest-600 flex-shrink-0 flex items-center justify-center"
                  >
                    <svg className="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                    </svg>
                  </button>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-farm-600 line-through">{todo.title}</p>
                    {todo.assigned_attendee_name && (
                      <p className="text-sm text-farm-400 mt-1">
                        {todo.assigned_attendee_name}
                      </p>
                    )}
                  </div>
                  <button
                    onClick={() => onDeleteTodo(todo.id)}
                    className="p-1 text-farm-300 hover:text-rust-600 flex-shrink-0"
                    title="Delete task"
                  >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
