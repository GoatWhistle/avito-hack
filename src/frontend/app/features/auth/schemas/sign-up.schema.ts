import { z } from 'zod'

export const SignUpSchema = z.object({
  email: z.email(),
  password: z.string().min(8),
  full_name: z.string().min(1).max(100),
})
