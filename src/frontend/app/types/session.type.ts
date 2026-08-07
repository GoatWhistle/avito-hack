import type { User, UserResponse } from './user.type'

export interface Session {
  token: string
  user: User
}

export interface SessionResponse {
  token: string
  user: UserResponse
}
