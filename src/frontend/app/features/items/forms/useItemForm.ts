import { useForm } from '@tanstack/react-form'
import { ItemFormSchema, type ItemFormValues } from '#/features/items/schemas'
import { kopeksToRubles, rublesToKopeks } from '#/features/items/lib'
import type { Item } from '#/features/items/types'

export interface ItemFormPayload {
  title: string
  description: string
  price: number
}

interface UseItemFormParams {
  item?: Item
  onSave: (payload: ItemFormPayload) => Promise<unknown>
}

export const toFormValues = (item?: Item): ItemFormValues => ({
  title: item?.title ?? '',
  description: item?.description ?? '',
  price: item ? kopeksToRubles(item.price) : 0,
})

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
      }),
  })
