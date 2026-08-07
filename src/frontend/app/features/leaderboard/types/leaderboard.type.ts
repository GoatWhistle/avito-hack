import { z } from 'zod'

export const leaderboardEntrySchema = z.object({
  user_id: z.string(),
  name: z.string(),
  level: z.number(),
  xp: z.number(),
  streak_days: z.number(),
  rank: z.number(),
})

export const leaderboardResponseSchema = z.object({
  items: z
    .array(leaderboardEntrySchema)
    .nullish()
    .transform((value) => value ?? []),
  my_rank: z.number().nullish(),
  next_cursor: z.string().optional().default(''),
})

export type LeaderboardEntry = z.infer<typeof leaderboardEntrySchema>
export type LeaderboardResponse = z.infer<typeof leaderboardResponseSchema>

export type LeaderboardQuery = {
  cursor?: string
  limit?: number
  around?: boolean
}

export type LeaderboardPage = {
  items: LeaderboardEntry[]
  myRank: number | null
  nextCursor: string
}
