export type UserRole = 'user' | 'moderator' | 'admin';

export interface User {
  id: string;
  email: string;
  fullName: string;
  role: UserRole;
  createdAt: string;
}

export interface Session {
  token: string;
  expiresAt: string;
  user: User;
}

export interface UserDto {
  id: string;
  email: string;
  full_name?: string;
  display_name?: string;
  role: UserRole;
  created_at: string;
}

export interface SessionDto {
  token: string;
  expires_at: string;
  user: UserDto;
}

export function toUser(dto: UserDto): User {
  return {
    id: dto.id,
    email: dto.email,
    fullName: dto.full_name ?? dto.display_name ?? '',
    role: dto.role,
    createdAt: dto.created_at,
  };
}

export function toSession(dto: SessionDto): Session {
  return {
    token: dto.token,
    expiresAt: dto.expires_at,
    user: toUser(dto.user),
  };
}

export function isModerator(user: User | null): boolean {
  return user?.role === 'moderator' || user?.role === 'admin';
}
