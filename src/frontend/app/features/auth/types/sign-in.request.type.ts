import type { SignInSchema } from '#/features/auth/schemas'
import { z } from 'zod'

export type SignInRequest = z.infer<typeof SignInSchema>
