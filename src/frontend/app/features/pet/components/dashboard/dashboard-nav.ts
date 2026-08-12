import {
  Gamepad2,
  Gift,
  Home,
  Medal,
  TrendingUp,
  Trophy,
  type LucideIcon,
} from 'lucide-react'

export type DashboardNavKey =
  | 'home'
  | 'progress'
  | 'games'
  | 'rewards'
  | 'achievements'
  | 'leaderboard'

export interface DashboardNavItem {
  key: DashboardNavKey
  icon: LucideIcon
  path: string
  end: boolean
}

export const dashboardNavItems: DashboardNavItem[] = [
  { key: 'home', icon: Home, path: '/pet', end: true },
  { key: 'progress', icon: TrendingUp, path: '/pet/progress', end: false },
  { key: 'games', icon: Gamepad2, path: '/pet/games', end: false },
  { key: 'rewards', icon: Gift, path: '/pet/rewards', end: false },
  { key: 'achievements', icon: Medal, path: '/pet/achievements', end: false },
  { key: 'leaderboard', icon: Trophy, path: '/pet/leaderboard', end: false },
]
