import { z } from 'zod'

export const petStageSchema = z.enum(['egg', 'baby', 'teen', 'adult', 'legend'])
export const petStateSchema = z.enum(['happy', 'neutral', 'sad', 'sleeping'])

export type PetStageValue = z.infer<typeof petStageSchema>
export type PetStateValue = z.infer<typeof petStateSchema>

export const petSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  name: z.string(),
  stage: petStageSchema,
  state: petStateSchema,
  level: z.number(),
  xp: z.number(),
  next_level_xp: z.number(),
  satiety: z.number(),
  happiness: z.number(),
  energy: z.number(),
  streak_days: z.number(),
  freezes: z.number(),
  is_hatched: z.boolean(),
  hatched_at: z.string().nullish(),
  last_checkin_date: z.string().nullish(),
  last_decay_time: z.string().optional(),
  updated_at: z.string().optional(),
})

export type Pet = z.infer<typeof petSchema>

export const progressSchema = z.object({
  level: z.number(),
  xp: z.number(),
  next_level_xp: z.number(),
  xp_to_next_level: z.number(),
  is_max_level: z.boolean(),
  stage: petStageSchema,
  streak_days: z.number(),
  freezes: z.number(),
})

export type Progress = z.infer<typeof progressSchema>

export const streakResultSchema = z.object({
  days: z.number(),
  continued: z.boolean(),
  freeze_used: z.boolean(),
  reset: z.boolean(),
  milestone_bonus: z.number(),
  milestone_reached: z.number(),
  freezes_left: z.number(),
})

export type StreakResult = z.infer<typeof streakResultSchema>

export const checkInResultSchema = z.object({
  pet: petSchema,
  xp_granted: z.number(),
  level: z.number(),
  previous_level: z.number(),
  next_level_xp: z.number(),
  unlocked_rewards: z.array(z.string()),
  streak: streakResultSchema,
})

export type CheckInResult = z.infer<typeof checkInResultSchema>
