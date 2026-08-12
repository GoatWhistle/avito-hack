import { z } from 'zod'
import {
  DESCRIPTION_MAX,
  MAX_PRICE_RUBLES,
  TITLE_MAX,
  TITLE_MIN,
} from '#/features/items/lib/constants'
import { itemCategories, itemConditions } from '#/features/items/types'

export const CUSTOM_CATEGORY_VALUE = 'other'
export const CUSTOM_CATEGORY_MAX = 60

export const ItemFormSchema = z
  .object({
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
    category: z.enum(itemCategories, { message: 'validation:required' }),
    customCategory: z
      .string()
      .trim()
      .max(CUSTOM_CATEGORY_MAX, { message: 'validation:maxLength' }),
    condition: z.enum(itemConditions, { message: 'validation:required' }),
  })
  .refine(
    (value) =>
      value.category !== CUSTOM_CATEGORY_VALUE ||
      value.customCategory.length > 0,
    { path: ['customCategory'], message: 'validation:required' },
  )

export type ItemFormValues = z.infer<typeof ItemFormSchema>
