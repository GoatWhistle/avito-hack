import { z } from 'zod'
import { emailSchema, signInPasswordSchema } from './credentials.schema'

export const SignInSchema = z.object({
  email: emailSchema,
  password: signInPasswordSchema,
})
