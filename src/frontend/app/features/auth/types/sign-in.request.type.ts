import type { SignInSchema } from '#/features/auth/schemas/sign-in.schema'
import { z } from 'zod'

export type SignInRequest = z.infer<typeof SignInSchema>
