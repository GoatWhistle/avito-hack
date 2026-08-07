import { z } from 'zod'

export const rewardKindSchema = z.enum(['promo', 'utility', 'cosmetic'])

export const rewardConditionTypeSchema = z.enum([
  'level',
  'streak',
  'achievement',
])

export const rewardStatusSchema = z.enum(['granted', 'activated', 'expired'])

export const rewardCatalogItemSchema = z.object({
  id: z.string(),
  title: z.string(),
  description: z.string(),
  kind: z.string(),
  condition_type: z.string(),
  condition_value: z.number(),
  unlocked: z.boolean(),
  claimed: z.boolean(),
  status: z.string().optional().default(''),
  progress_current: z.number(),
  progress_target: z.number(),
})

export const myRewardItemSchema = z.object({
  reward_id: z.string(),
  title: z.string(),
  description: z.string(),
  kind: z.string(),
  status: z.string(),
  code: z.string().optional().default(''),
  granted_at: z.string(),
  activated_at: z.string().nullish(),
  expires_at: z.string().nullish(),
})

export const activateRewardResponseSchema = z.object({
  reward_id: z.string(),
  code: z.string(),
  status: z.string(),
})

export const badgeItemSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string(),
  icon_url: z.string().optional().default(''),
  earned_at: z.string().optional().default(''),
})

export const listResponseSchema = <T extends z.ZodTypeAny>(item: T) =>
  z.object({
    items: z
      .array(item)
      .nullish()
      .transform((value) => value ?? []),
    next_cursor: z.string().optional().default(''),
  })

export type RewardKind = z.infer<typeof rewardKindSchema>
export type RewardConditionType = z.infer<typeof rewardConditionTypeSchema>
export type RewardStatus = z.infer<typeof rewardStatusSchema>
export type RewardCatalogItem = z.infer<typeof rewardCatalogItemSchema>
export type MyRewardItem = z.infer<typeof myRewardItemSchema>
export type ActivateRewardResponse = z.infer<
  typeof activateRewardResponseSchema
>
export type BadgeItem = z.infer<typeof badgeItemSchema>
