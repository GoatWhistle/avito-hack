import { z } from 'zod'

export const emailSchema = z
  .string()
  .trim()
  .min(1, { message: 'validation:required' })
  .refine((value) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value), {
    message: 'validation:email',
  })

export const passwordSchema = z
  .string()
  .min(1, { message: 'validation:required' })

export const fullNameSchema = z
  .string()
  .trim()
  .min(1, { message: 'validation:required' })

export const credentialsShape = {
  email: emailSchema,
  password: passwordSchema,
}
