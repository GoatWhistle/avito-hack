import { useForm } from '@tanstack/react-form'
import {
  CUSTOM_CATEGORY_VALUE,
  ItemFormSchema,
  type ItemFormValues,
} from '#/features/items/schemas'
import { kopeksToRubles, rublesToKopeks } from '#/features/items/lib'
import { itemCategories, type Item } from '#/features/items/types'

export interface ItemFormPayload {
  title: string
  description: string
  price: number
  attributes: Record<string, string>
}

interface UseItemFormParams {
  item?: Item
  onSave: (payload: ItemFormPayload) => Promise<unknown>
}

const isKnownCategory = (value: string): value is (typeof itemCategories)[number] =>
  (itemCategories as readonly string[]).includes(value)

export const toFormValues = (item?: Item): ItemFormValues => {
  const rawCategory = item?.attributes?.category ?? ''
  const category = isKnownCategory(rawCategory)
    ? rawCategory
    : rawCategory
      ? CUSTOM_CATEGORY_VALUE
      : itemCategories[0]
  const customCategory =
    rawCategory && !isKnownCategory(rawCategory) ? rawCategory : ''

  return {
    title: item?.title ?? '',
    description: item?.description ?? '',
    price: item ? kopeksToRubles(item.price) : 0,
    category,
    customCategory,
    condition: item?.attributes?.condition === 'new' ? 'new' : 'used',
  }
}

export const useItemForm = ({ item, onSave }: UseItemFormParams) =>
  useForm({
    formId: item ? `item-edit-${item.id}` : 'item-create',
    defaultValues: toFormValues(item),
    validators: {
      onChange: ItemFormSchema,
    },
    onSubmit: ({ value }) =>
      onSave({
        title: value.title.trim(),
        description: value.description.trim(),
        price: rublesToKopeks(value.price),
        attributes: {
          category:
            value.category === CUSTOM_CATEGORY_VALUE
              ? value.customCategory.trim()
              : value.category,
          condition: value.condition,
        },
      }),
  })
