import { z } from 'zod'

export const questItemSchema = z.object({
  id: z.string(),
  action: z.string(),
  target: z.number(),
  reward_xp: z.number(),
  progress_current: z.number(),
  completed: z.boolean(),
  claimed: z.boolean(),
})

export const questListResponseSchema = z.object({
  items: z
    .array(questItemSchema)
    .nullish()
    .transform((value) => value ?? []),
  next_cursor: z.string().optional().default(''),
})

export type QuestItem = z.infer<typeof questItemSchema>

export interface QuestProgressView {
  quest: QuestItem
  current: number
  target: number
  remaining: number
  percent: number
  completed: boolean
  claimed: boolean
}

export interface QuestSummary {
  completed: number
  total: number
  earnedXP: number
  pendingXP: number
}
