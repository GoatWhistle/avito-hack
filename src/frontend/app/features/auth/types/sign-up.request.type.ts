import type { SignUpSchema } from '#/features/auth/schemas'
import { z } from 'zod'

export type SignUpRequest = z.infer<typeof SignUpSchema>
