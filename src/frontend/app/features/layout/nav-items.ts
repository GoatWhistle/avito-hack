import {
  Heart,
  ListOrdered,
  PawPrint,
  Package,
  Store,
  Trophy,
  type LucideIcon,
} from 'lucide-react'

export interface NavItem {
  key: 'pet' | 'items' | 'myItems' | 'favorites' | 'rewards' | 'leaderboard'
  to: string
  icon: LucideIcon
  primary: boolean
}

export const navItems: NavItem[] = [
  { key: 'pet', to: '/pet', icon: PawPrint, primary: true },
  { key: 'items', to: '/items', icon: Store, primary: true },
  { key: 'myItems', to: '/items/mine', icon: Package, primary: false },
  { key: 'favorites', to: '/favorites', icon: Heart, primary: true },
  { key: 'rewards', to: '/rewards', icon: Trophy, primary: true },
  { key: 'leaderboard', to: '/leaderboard', icon: ListOrdered, primary: false },
]
