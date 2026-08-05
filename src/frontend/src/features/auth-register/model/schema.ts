import type { TFunction } from 'i18next';
import { z } from 'zod';

export function buildRegisterSchema(t: TFunction<'validation'>) {
  return z.object({
    email: z.string().min(1, t('required')).email(t('email')),
    password: z.string().min(1, t('required')),
    fullName: z.string().trim().min(1, t('required')),
  });
}

export type RegisterFormValues = z.infer<ReturnType<typeof buildRegisterSchema>>;
