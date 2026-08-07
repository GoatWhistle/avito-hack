import { z } from 'zod'
import { credentialsShape, fullNameSchema } from './credentials.schema'

export const SignUpSchema = z.object({
  ...credentialsShape,
  fullName: fullNameSchema,
})
