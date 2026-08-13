import { z } from 'zod'

export const lotterySymbols = [
  'bicycle',
  'smartphone',
  'sofa',
  'sneakers',
  'delivery',
  'promotion',
] as const

export const lotterySymbolSchema = z.enum(lotterySymbols)
export type LotterySymbol = z.infer<typeof lotterySymbolSchema>

export const lotteryRunStateSchema = z.enum(['active', 'won', 'lost'])
export type LotteryRunState = z.infer<typeof lotteryRunStateSchema>

const slotIndexSchema = z.number().int().min(0).max(8)

export const lotterySlotSchema = z.discriminatedUnion('opened', [
  z.object({
    index: slotIndexSchema,
    opened: z.literal(false),
    symbol: z.undefined().optional(),
  }),
  z.object({
    index: slotIndexSchema,
    opened: z.literal(true),
    symbol: lotterySymbolSchema,
  }),
])

export const lotteryPrizeSchema = z.object({
  id: z.string(),
  title: z.string(),
  description: z.string(),
  benefit_type: z.string(),
  benefit_value: z.number().int(),
  scope_type: z.string(),
  scope_value: z.string(),
  code: z.string(),
  expires_at: z.string(),
})

export const lotteryPrizeCatalogItemSchema = lotteryPrizeSchema
  .omit({ code: true, expires_at: true })
  .extend({ symbol: lotterySymbolSchema })

export const lotteryRunSchema = z
  .object({
    id: z.string(),
    state: lotteryRunStateSchema,
    slots: z.array(lotterySlotSchema).length(9),
    prize: lotteryPrizeSchema.optional(),
    created_at: z.string(),
  })
  .superRefine((run, context) => {
    const indexes = new Set(run.slots.map((slot) => slot.index))
    if (indexes.size !== 9) {
      context.addIssue({
        code: 'custom',
        path: ['slots'],
        message: 'lottery slots must have unique indexes',
      })
    }

    if ((run.state === 'won') !== Boolean(run.prize)) {
      context.addIssue({
        code: 'custom',
        path: ['prize'],
        message: 'only a won lottery run may contain a prize',
      })
    }

    if (run.state === 'lost' && run.slots.some((slot) => !slot.opened)) {
      context.addIssue({
        code: 'custom',
        path: ['slots'],
        message: 'a lost lottery run must have every slot opened',
      })
    }
  })

export const lotteryStateSchema = z
  .object({
    available: z.boolean(),
    week_start: z.string(),
    next_available_at: z.string(),
    run: lotteryRunSchema.optional(),
  })
  .superRefine((state, context) => {
    const expectedAvailable = !state.run || state.run.state === 'active'
    if (state.available !== expectedAvailable) {
      context.addIssue({
        code: 'custom',
        path: ['available'],
        message: 'lottery availability does not match its run state',
      })
    }
  })

export const lotteryRevealSchema = z
  .object({
    index: z.number().int().min(0).max(8),
    symbol: lotterySymbolSchema,
    run: lotteryRunSchema,
  })
  .superRefine((reveal, context) => {
    const slot = reveal.run.slots.find(({ index }) => index === reveal.index)
    if (!slot?.opened || slot.symbol !== reveal.symbol) {
      context.addIssue({
        code: 'custom',
        path: ['run', 'slots'],
        message: 'revealed slot does not match the returned run',
      })
    }
  })

export type LotterySlot = z.infer<typeof lotterySlotSchema>
export type LotteryPrize = z.infer<typeof lotteryPrizeSchema>
export type LotteryPrizeCatalogItem = z.infer<
  typeof lotteryPrizeCatalogItemSchema
>
export type LotteryRun = z.infer<typeof lotteryRunSchema>
export type LotteryState = z.infer<typeof lotteryStateSchema>
export type LotteryReveal = z.infer<typeof lotteryRevealSchema>
