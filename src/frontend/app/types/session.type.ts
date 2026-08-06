import type { User } from './user.type'

export interface Session {
  token: string
  user: User
}
