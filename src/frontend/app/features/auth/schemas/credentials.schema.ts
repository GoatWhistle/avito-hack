import { z } from 'zod'

export const PASSWORD_MIN_LENGTH = 8
export const PASSWORD_MAX_LENGTH = 72

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
  .min(PASSWORD_MIN_LENGTH, { message: 'validation:passwordTooShort' })
  .max(PASSWORD_MAX_LENGTH, { message: 'validation:passwordTooLong' })

export const signInPasswordSchema = z
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
