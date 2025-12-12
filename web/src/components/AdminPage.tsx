import { useState, useEffect } from 'react'
import { api } from '../api/client'
import type { User } from '../api/types'

interface AdminPageProps {
  onBack: () => void
}

export function AdminPage({ onBack }: AdminPageProps) {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    loadUsers()
  }, [])

  const loadUsers = async () => {
    try {
      setLoading(true)
      const data = await api.listUsers()
      setUsers(data || [])
      setError('')
    } catch (err) {
      setError('Failed to load users. You may not have admin access.')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleTogglePermission = async (userId: string, currentValue: boolean) => {
    try {
      await api.updateUserPermissions(userId, { can_create_events: !currentValue })
      await loadUsers()
    } catch (err) {
      setError('Failed to update user permissions')
      console.error(err)
    }
  }

  return (
    <div>
      <button
        onClick={onBack}
        className="mb-6 text-forest-600 hover:text-forest-800 flex items-center gap-2"
      >
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
        </svg>
        Back to Events
      </button>

      <div className="bg-white rounded-lg shadow-md border border-farm-200 p-6">
        <h2 className="text-xl font-semibold text-farm-800 mb-6">User Management</h2>

        {error && (
          <div className="mb-6 p-4 bg-rust-100 text-rust-800 rounded-lg border border-rust-300">
            {error}
          </div>
        )}

        {loading ? (
          <div className="flex justify-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-4 border-forest-600 border-t-transparent"></div>
          </div>
        ) : (
          <div className="space-y-4">
            {users.length === 0 ? (
              <p className="text-farm-500 text-center py-4">No users found</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-farm-200">
                      <th className="text-left py-3 px-2 text-farm-600 font-medium">User</th>
                      <th className="text-left py-3 px-2 text-farm-600 font-medium">Role</th>
                      <th className="text-center py-3 px-2 text-farm-600 font-medium">Can Create Events</th>
                    </tr>
                  </thead>
                  <tbody>
                    {users.map((user) => (
                      <tr key={user.id} className="border-b border-farm-100 hover:bg-cream-50">
                        <td className="py-3 px-2">
                          <div className="flex items-center gap-3">
                            <img
                              src={user.picture}
                              alt={user.name}
                              className="w-8 h-8 rounded-full"
                              referrerPolicy="no-referrer"
                            />
                            <div>
                              <div className="font-medium text-farm-800">{user.name}</div>
                              <div className="text-sm text-farm-500">{user.email}</div>
                            </div>
                          </div>
                        </td>
                        <td className="py-3 px-2">
                          {user.is_admin ? (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-forest-100 text-forest-800">
                              Admin
                            </span>
                          ) : (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-farm-100 text-farm-600">
                              User
                            </span>
                          )}
                        </td>
                        <td className="py-3 px-2 text-center">
                          {user.is_admin ? (
                            <span className="text-farm-400 text-sm">Always</span>
                          ) : (
                            <button
                              onClick={() => handleTogglePermission(user.id, user.can_create_events)}
                              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                                user.can_create_events ? 'bg-forest-600' : 'bg-farm-300'
                              }`}
                            >
                              <span
                                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                                  user.can_create_events ? 'translate-x-6' : 'translate-x-1'
                                }`}
                              />
                            </button>
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        <div className="mt-6 pt-4 border-t border-farm-200">
          <p className="text-sm text-farm-500">
            Users with "Can Create Events" permission can create new events. Admins always have this permission.
          </p>
        </div>
      </div>
    </div>
  )
}
