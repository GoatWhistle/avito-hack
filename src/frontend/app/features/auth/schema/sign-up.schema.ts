import { z } from 'zod'

export const SignUpSchema = z.object({
  email: z.email(),
  password: z.string().min(8),
  fullName: z.string().min(50),
})
