import { Apple, Smile, Zap, type LucideIcon } from 'lucide-react'
import type { Pet } from '#/features/pet/types'

export type DashboardStatKey =
  'energy' | 'happiness' | 'satiety' | 'streak' | 'xp'

export interface DashboardStat {
  key: DashboardStatKey
  icon: LucideIcon
  value: string | number
  percent: number
}

export interface DashboardStatStyle {
  icon: string
  indicator: string
}

export const STREAK_TARGET_DAYS = 30

export const DASHBOARD_STAT_STYLES: Record<
  DashboardStatKey,
  DashboardStatStyle
> = {
  energy: {
    icon: 'bg-stat-energy-subtle text-stat-energy-subtle-foreground',
    indicator: 'bg-stat-energy-fill',
  },
  happiness: {
    icon: 'bg-stat-happiness-subtle text-stat-happiness-subtle-foreground',
    indicator: 'bg-stat-happiness-fill',
  },
  satiety: {
    icon: 'bg-stat-satiety-subtle text-stat-satiety-subtle-foreground',
    indicator: 'bg-stat-satiety-fill',
  },
  streak: {
    icon: 'bg-stat-streak-subtle text-stat-streak-subtle-foreground',
    indicator: 'bg-stat-streak-fill',
  },
  xp: {
    icon: 'bg-stat-xp-subtle text-stat-xp-subtle-foreground',
    indicator: 'bg-stat-xp-fill',
  },
}

const clamp = (value: number): number => Math.max(0, Math.min(100, value))

export const dashboardStats = (pet: Pet): DashboardStat[] => {
  return [
    {
      key: 'energy',
      icon: Zap,
      value: pet.energy,
      percent: clamp(pet.energy),
    },
    {
      key: 'happiness',
      icon: Smile,
      value: pet.happiness,
      percent: clamp(pet.happiness),
    },
    {
      key: 'satiety',
      icon: Apple,
      value: pet.satiety,
      percent: clamp(pet.satiety),
    },
  ]
}
