import { z } from 'zod'
import type { SignInSchema } from '#/features/auth/schemas'

export type SignInRequest = z.infer<typeof SignInSchema>
