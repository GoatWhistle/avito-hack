import type { SignUpSchema } from '#/features/auth/schemas/sign-up.schema'
import { z } from 'zod'

export type SignUpRequest = z.infer<typeof SignUpSchema>
