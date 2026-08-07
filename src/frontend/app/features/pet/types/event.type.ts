import { z } from 'zod'
import { petSchema } from './pet.type'

export const xpGainedPayloadSchema = z.object({
  amount: z.number(),
  reason: z.string().optional(),
  total: z.number().optional(),
})

export const levelUpPayloadSchema = z.object({
  level: z.number(),
})

export const rewardGrantedPayloadSchema = z.object({
  reward_id: z.string(),
  title: z.string(),
})

export const streakUpdatedPayloadSchema = z.object({
  days: z.number(),
  milestone: z.boolean().optional(),
})

export const petEventSchema = z.discriminatedUnion('type', [
  z.object({ type: z.literal('pet.state'), payload: petSchema }),
  z.object({ type: z.literal('pet.updated'), payload: petSchema }),
  z.object({ type: z.literal('pet.hatched'), payload: petSchema }),
  z.object({ type: z.literal('xp.gained'), payload: xpGainedPayloadSchema }),
  z.object({ type: z.literal('level.up'), payload: levelUpPayloadSchema }),
  z.object({
    type: z.literal('reward.granted'),
    payload: rewardGrantedPayloadSchema,
  }),
  z.object({
    type: z.literal('streak.updated'),
    payload: streakUpdatedPayloadSchema,
  }),
])

export type PetEvent = z.infer<typeof petEventSchema>
export type PetEventType = PetEvent['type']

export const parsePetEvent = (raw: unknown): PetEvent | null => {
  const parsed = petEventSchema.safeParse(raw)
  return parsed.success ? parsed.data : null
}
