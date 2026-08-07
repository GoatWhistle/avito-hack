export interface BadgeResponse {
  id: string
  name: string
  description: string
  icon_url: string
  earned_at?: string
}

export interface RaccoonProfileResponse {
  id: string
  user_id: string
  name: string
  level: number
  xp: number
  xp_to_next_level: number
  current_streak: number
  badges?: BadgeResponse[]
}

export interface ProgressBadge {
  id: string
  name: string
  description: string
  iconUrl: string
  earnedAt: string | null
}

export interface Progress {
  id: string
  name: string
  level: number
  xp: number
  xpToNextLevel: number
  currentStreak: number
  badges: ProgressBadge[]
  earnedBadgeCount: number
}

export const levelThresholds = [
  0, 5, 12, 22, 35, 52, 72, 95, 122, 155, 195, 240, 290, 350, 420,
] as const

export const levelFloor = (level: number): number =>
  levelThresholds[Math.min(Math.max(level, 1), levelThresholds.length) - 1] ?? 0

export const xpProgressRatio = (progress: Progress): number => {
  if (progress.xpToNextLevel <= 0) return 1

  const floor = levelFloor(progress.level)
  const span = progress.xpToNextLevel - floor
  if (span <= 0) return 1

  const gained = progress.xp - floor

  return Math.min(Math.max(gained / span, 0), 1)
}
