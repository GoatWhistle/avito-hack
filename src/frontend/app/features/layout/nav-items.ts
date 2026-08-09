import {
  Heart,
  Package,
  PawPrint,
  Store,
  type LucideIcon,
} from 'lucide-react'

export interface NavItem {
  key: 'pet' | 'items' | 'myItems' | 'favorites'
  to: string
  icon: LucideIcon
  primary: boolean
}

export const navItems: NavItem[] = [
  { key: 'pet', to: '/pet', icon: PawPrint, primary: true },
  { key: 'items', to: '/items', icon: Store, primary: true },
]

export const accountNavItems: NavItem[] = [
  { key: 'myItems', to: '/items/mine', icon: Package, primary: false },
  { key: 'favorites', to: '/favorites', icon: Heart, primary: false },
]
