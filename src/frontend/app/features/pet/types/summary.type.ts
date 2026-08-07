import { z } from 'zod'
import { petStageSchema, petStateSchema } from './pet.type'

export const adviceActionSchema = z.enum([
  'add_photo',
  'expand_description',
  'refresh_listing',
  'set_price',
  'check_in',
])

export type AdviceAction = z.infer<typeof adviceActionSchema>

export const adviceSchema = z.object({
  text: z.string(),
  item_id: z.string().nullish(),
  action: z.string(),
})

export type Advice = z.infer<typeof adviceSchema>

export const summaryActionSchema = z.object({
  action: z.string(),
  count: z.number(),
  amount: z.number(),
})

export const summaryFactsSchema = z.object({
  total_xp: z.number(),
  actions: z.array(summaryActionSchema),
  level: z.number(),
  previous_level: z.number(),
  leveled_up: z.boolean(),
  xp: z.number(),
  next_level_xp: z.number(),
  xp_to_next_level: z.number(),
  rewards: z.array(z.string()),
  badges: z.array(z.string()),
  stage: petStageSchema,
  state: petStateSchema,
  satiety: z.number(),
  happiness: z.number(),
  energy: z.number(),
  streak_days: z.number(),
  streak_broken: z.boolean(),
  leaderboard_rank: z.number().nullish(),
  issues_count: z.number(),
})

export type SummaryFacts = z.infer<typeof summaryFactsSchema>

export const dailySummarySchema = z.object({
  id: z.string(),
  date: z.string(),
  message: z.string(),
  advice: adviceSchema.nullish(),
  generated_by: z.string(),
  facts: summaryFactsSchema,
  created_at: z.string().optional(),
})

export type DailySummary = z.infer<typeof dailySummarySchema>
