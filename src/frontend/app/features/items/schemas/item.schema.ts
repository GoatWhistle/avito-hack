import { z } from 'zod'
import {
  DESCRIPTION_MAX,
  MAX_PRICE_RUBLES,
  TITLE_MAX,
  TITLE_MIN,
} from '#/features/items/lib/constants'

export const ItemFormSchema = z.object({
  title: z
    .string()
    .trim()
    .min(TITLE_MIN, { message: 'validation:titleTooShort' })
    .max(TITLE_MAX, { message: 'validation:maxLength' }),
  description: z
    .string()
    .trim()
    .max(DESCRIPTION_MAX, { message: 'validation:maxLength' }),
  price: z
    .number({ message: 'validation:invalidNumber' })
    .min(0, { message: 'validation:positiveNumber' })
    .max(MAX_PRICE_RUBLES, { message: 'validation:max' }),
})

export type ItemFormValues = z.infer<typeof ItemFormSchema>
