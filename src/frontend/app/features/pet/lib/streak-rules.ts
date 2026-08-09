import type { Pet } from '#/features/pet/types'

export const STREAK_BONUS_DAYS = 7
export const STREAK_BONUS_FACTOR = 1.5
export const STREAK_RESET_HOURS = 24
export const MAX_FREEZES = 3
export const STREAK_MILESTONES = [3, 7, 14, 30] as const

export const MILESTONE_BONUS_XP: Record<number, number> = {
  3: 30,
  7: 70,
  14: 150,
  30: 300,
}

export type StreakRuleKey = 'bonus' | 'reset' | 'freeze'
export type StreakRuleTone = 'active' | 'pending' | 'muted'

export interface StreakRuleView {
  key: StreakRuleKey
  icon: string
  tone: StreakRuleTone
  badgeCount: number
}

export const isBonusActive = (streakDays: number): boolean =>
  streakDays >= STREAK_BONUS_DAYS

export const daysToBonus = (streakDays: number): number =>
  Math.max(0, STREAK_BONUS_DAYS - Math.max(0, Math.trunc(streakDays)))

export const nextMilestone = (streakDays: number): number | null => {
  const safe = Math.max(0, Math.trunc(streakDays))
  const found = STREAK_MILESTONES.find((milestone) => milestone > safe)

  return found ?? null
}

export const daysToNextMilestone = (streakDays: number): number | null => {
  const milestone = nextMilestone(streakDays)

  return milestone === null ? null : milestone - Math.max(0, streakDays)
}

export const streakRules = (pet: Pet): StreakRuleView[] => {
  const streakDays = Math.max(0, Math.trunc(pet.streak_days))
  const freezes = Math.max(0, Math.trunc(pet.freezes))

  return [
    {
      key: 'bonus',
      icon: '🔥',
      tone: isBonusActive(streakDays) ? 'active' : 'pending',
      badgeCount: isBonusActive(streakDays) ? streakDays : daysToBonus(streakDays),
    },
    { key: 'reset', icon: '⏳', tone: 'muted', badgeCount: STREAK_RESET_HOURS },
    {
      key: 'freeze',
      icon: '🛟',
      tone: freezes > 0 ? 'active' : 'muted',
      badgeCount: freezes,
    },
  ]
}
