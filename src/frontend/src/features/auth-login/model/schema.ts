import type { TFunction } from 'i18next';
import { z } from 'zod';

const MIN_PASSWORD_LENGTH = 8;
const MAX_PASSWORD_LENGTH = 72;
const MAX_EMAIL_LENGTH = 254;

export function buildLoginSchema(t: TFunction<'validation'>) {
  return z.object({
    email: z
      .string()
      .min(1, t('required'))
      .max(MAX_EMAIL_LENGTH, t('maxLength', { count: MAX_EMAIL_LENGTH }))
      .email(t('email')),
    password: z
      .string()
      .min(MIN_PASSWORD_LENGTH, t('minLength', { count: MIN_PASSWORD_LENGTH }))
      .max(MAX_PASSWORD_LENGTH, t('maxLength', { count: MAX_PASSWORD_LENGTH })),
  });
}

export type LoginFormValues = z.infer<ReturnType<typeof buildLoginSchema>>;
