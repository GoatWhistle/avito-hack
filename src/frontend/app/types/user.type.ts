export interface User {
  id: string
  email: string
  fullName: string
  role: string
  createdAt: string
}

export interface UserResponse {
  id: string
  email: string
  full_name: string
  role: string
  created_at: string
}

export const toUser = (response: UserResponse): User => ({
  id: response.id,
  email: response.email,
  fullName: response.full_name,
  role: response.role,
  createdAt: response.created_at,
})
