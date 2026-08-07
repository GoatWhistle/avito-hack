import type { Pet } from '#/features/pet/types'

export const MAX_LEVEL = 15

export const STAT_CRITICAL_BELOW = 30
export const STAT_WARNING_BELOW = 60

export type StatKey = 'satiety' | 'happiness' | 'energy'
export type StatTone = 'critical' | 'warning' | 'good'

export const statTone = (value: number): StatTone => {
  if (value < STAT_CRITICAL_BELOW) return 'critical'
  if (value < STAT_WARNING_BELOW) return 'warning'

  return 'good'
}

export const clampPercent = (value: number): number =>
  Math.max(0, Math.min(100, Math.round(value)))

export interface StatView {
  key: StatKey
  value: number
  tone: StatTone
}

export const statViews = (pet: Pet): StatView[] =>
  (['satiety', 'happiness', 'energy'] as StatKey[]).map((key) => ({
    key,
    value: clampPercent(pet[key]),
    tone: statTone(pet[key]),
  }))

export interface LevelProgress {
  level: number
  xp: number
  nextLevelXp: number
  xpToNext: number
  percent: number
  isMaxLevel: boolean
}

export const levelProgress = (pet: Pet): LevelProgress => {
  const isMaxLevel = pet.level >= MAX_LEVEL || pet.next_level_xp <= 0
  const xpToNext = isMaxLevel ? 0 : Math.max(0, pet.next_level_xp - pet.xp)
  const percent = isMaxLevel
    ? 100
    : clampPercent((pet.xp / Math.max(1, pet.next_level_xp)) * 100)

  return {
    level: pet.level,
    xp: pet.xp,
    nextLevelXp: pet.next_level_xp,
    xpToNext,
    percent,
    isMaxLevel,
  }
}

export type NextStepKey =
  | 'checkIn'
  | 'feedHungry'
  | 'cheerUp'
  | 'rest'
  | 'hatch'
  | 'keepGoing'
  | 'maxLevel'

export interface NextStepInput {
  pet: Pet
  canCheckIn: boolean
}

export const nextStep = ({ pet, canCheckIn }: NextStepInput): NextStepKey => {
  if (!pet.is_hatched) return 'hatch'
  if (canCheckIn) return 'checkIn'
  if (pet.satiety < STAT_CRITICAL_BELOW) return 'feedHungry'
  if (pet.energy < STAT_CRITICAL_BELOW) return 'rest'
  if (pet.happiness < STAT_WARNING_BELOW) return 'cheerUp'
  if (pet.level >= MAX_LEVEL) return 'maxLevel'

  return 'keepGoing'
}

export const isSameDay = (
  iso: string | null | undefined,
  now: Date,
): boolean => {
  if (iso === null || iso === undefined || iso === '') return false
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return false

  return date.toDateString() === now.toDateString()
}

export const canCheckInToday = (pet: Pet, now: Date = new Date()): boolean =>
  !isSameDay(pet.last_checkin_date, now)
