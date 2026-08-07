import { type z } from 'zod'
import type { SignUpSchema } from '#/features/auth/schemas'

export type SignUpRequest = z.infer<typeof SignUpSchema>
