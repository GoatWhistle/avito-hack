import { z } from 'zod';

import { PET_STAGES } from './types';

export const petStageSchema = z.enum(PET_STAGES);

export const petStateDtoSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  name: z.string(),
  stage: petStageSchema,
  level: z.number(),
  xp: z.number(),
  next_level_xp: z.number(),
  satiety: z.number(),
  happiness: z.number(),
  streak_days: z.number(),
  last_checkin_date: z.string().nullish(),
  last_decay_time: z.string(),
  updated_at: z.string(),
});

export const badgeDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().default(''),
  icon_url: z.string().default(''),
  earned_at: z.string().nullish(),
});

export const raccoonProfileDtoSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  name: z.string(),
  level: z.number(),
  xp: z.number(),
  xp_to_next_level: z.number(),
  current_streak: z.number(),
  badges: z.array(badgeDtoSchema).default([]),
});

export const badgeListSchema = z.array(badgeDtoSchema);

export const claimRewardDtoSchema = z.object({
  reward_id: z.string(),
  promocode: z.string(),
});

export const xpGainedPayloadSchema = z.object({
  amount: z.number(),
  reason: z.string().default(''),
  total: z.number().default(0),
});

export const levelUpPayloadSchema = z.object({
  level: z.number(),
  unlocked_rewards: z.array(z.string()).default([]),
});

export const rewardGrantedPayloadSchema = z.object({
  reward_id: z.string(),
  title: z.string().default(''),
});

export const streakUpdatedPayloadSchema = z.object({
  days: z.number(),
  milestone: z.boolean().default(false),
});

export const serverMessageSchema = z.object({
  type: z.string(),
  request_id: z.string().optional(),
  payload: z.unknown().optional(),
});

export type PetStateDto = z.infer<typeof petStateDtoSchema>;
export type BadgeDto = z.infer<typeof badgeDtoSchema>;
export type RaccoonProfileDto = z.infer<typeof raccoonProfileDtoSchema>;
export type ClaimRewardDto = z.infer<typeof claimRewardDtoSchema>;
export type ServerMessage = z.infer<typeof serverMessageSchema>;
