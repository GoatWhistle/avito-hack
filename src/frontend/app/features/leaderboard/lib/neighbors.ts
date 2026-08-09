import type { LeaderboardEntry } from '#/features/leaderboard/types'

export const NEIGHBOR_RADIUS = 3

export const isMe = (entry: LeaderboardEntry, myRank: number | null) =>
  myRank !== null && entry.rank === myRank

export const neighborsOf = (
  items: LeaderboardEntry[],
  myRank: number | null,
  radius = NEIGHBOR_RADIUS,
): LeaderboardEntry[] => {
  if (myRank === null) return []

  return items
    .filter((entry) => Math.abs(entry.rank - myRank) <= radius)
    .sort((a, b) => a.rank - b.rank)
}

export const chaseTarget = (
  items: LeaderboardEntry[],
  myRank: number | null,
): { target: LeaderboardEntry; gap: number } | null => {
  if (myRank === null || myRank <= 1) return null

  const me = items.find((entry) => entry.rank === myRank)
  const ahead = items.find((entry) => entry.rank === myRank - 1)

  if (!me || !ahead) return null

  return { target: ahead, gap: Math.max(0, ahead.xp - me.xp) }
}

export const mergePages = (
  pages: Array<{ items: LeaderboardEntry[] }>,
): LeaderboardEntry[] => {
  const seen = new Map<string, LeaderboardEntry>()

  for (const page of pages) {
    for (const entry of page.items) {
      if (!seen.has(entry.user_id)) seen.set(entry.user_id, entry)
    }
  }

  return [...seen.values()].sort((a, b) => a.rank - b.rank)
}

export const medalOf = (rank: number): 'gold' | 'silver' | 'bronze' | null => {
  if (rank === 1) return 'gold'
  if (rank === 2) return 'silver'
  if (rank === 3) return 'bronze'

  return null
}

export const initialsOf = (name: string): string =>
  name
    .split(/\s+/)
    .filter(Boolean)
    .map((part) => [...part][0] ?? '')
    .join('')
    .toUpperCase()
    .slice(0, 2)
