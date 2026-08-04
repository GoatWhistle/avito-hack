import type { TFunction } from 'i18next';
import { z } from 'zod';

const MIN_TITLE_LENGTH = 3;
const MAX_TITLE_LENGTH = 200;
const MAX_DESCRIPTION_LENGTH = 5000;

export function buildItemSchema(t: TFunction<'validation'>) {
  return z.object({
    title: z
      .string()
      .min(MIN_TITLE_LENGTH, t('minLength', { count: MIN_TITLE_LENGTH }))
      .max(MAX_TITLE_LENGTH, t('maxLength', { count: MAX_TITLE_LENGTH })),
    description: z.string().max(MAX_DESCRIPTION_LENGTH, t('maxLength', { count: MAX_DESCRIPTION_LENGTH })),
    price: z.number({ invalid_type_error: t('integer') }).nonnegative(t('nonNegative')),
  });
}

export type ItemFormValues = z.infer<ReturnType<typeof buildItemSchema>>;
