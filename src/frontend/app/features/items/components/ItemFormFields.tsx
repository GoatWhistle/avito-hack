import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
  Input,
  Select,
} from '#/components/ui'
import { DESCRIPTION_MAX, TITLE_MAX } from '#/features/items/lib'
import { useFieldErrors } from '#/features/items/forms'
import { CUSTOM_CATEGORY_VALUE, CUSTOM_CATEGORY_MAX } from '#/features/items/schemas'
import { itemCategories, itemConditions } from '#/features/items/types'
import type { useItemForm } from '#/features/items/forms'

interface ItemFormFieldsProps {
  form: ReturnType<typeof useItemForm>
}

export function ItemFormFields({ form }: ItemFormFieldsProps) {
  const { t } = useTranslation('items')
  const toErrors = useFieldErrors()

  const categoryOptions = useMemo(
    () =>
      itemCategories.map((category) => ({
        value: category,
        label: t(`category.${category}`),
      })),
    [t],
  )

  const conditionOptions = useMemo(
    () =>
      itemConditions.map((condition) => ({
        value: condition,
        label: t(`condition.${condition}`),
      })),
    [t],
  )

  return (
    <FieldGroup>
      <form.Field
        name="title"
        children={(field) => (
          <Field data-invalid={field.state.meta.errors.length > 0}>
            <FieldLabel htmlFor={field.name}>{t('fields.title')}</FieldLabel>
            <Input
              id={field.name}
              name={field.name}
              maxLength={TITLE_MAX}
              value={field.state.value}
              placeholder={t('placeholders.title')}
              aria-invalid={field.state.meta.errors.length > 0}
              onBlur={field.handleBlur}
              onChange={(event) => field.handleChange(event.target.value)}
            />
            <FieldError errors={toErrors(field.state.meta.errors, 'title')} />
          </Field>
        )}
      />

      <form.Field
        name="description"
        children={(field) => (
          <Field data-invalid={field.state.meta.errors.length > 0}>
            <FieldLabel htmlFor={field.name}>
              {t('fields.description')}
            </FieldLabel>
            <textarea
              id={field.name}
              name={field.name}
              rows={6}
              maxLength={DESCRIPTION_MAX}
              value={field.state.value}
              placeholder={t('placeholders.description')}
              aria-invalid={field.state.meta.errors.length > 0}
              onBlur={field.handleBlur}
              onChange={(event) => field.handleChange(event.target.value)}
              className="w-full resize-y rounded-lg border border-input bg-transparent px-2.5 py-2 text-base outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 md:text-sm dark:bg-input/30"
            />
            <FieldError
              errors={toErrors(field.state.meta.errors, 'description')}
            />
          </Field>
        )}
      />

      <form.Field
        name="price"
        children={(field) => (
          <Field data-invalid={field.state.meta.errors.length > 0}>
            <FieldLabel htmlFor={field.name}>{t('fields.price')}</FieldLabel>
            <Input
              id={field.name}
              name={field.name}
              type="number"
              min={0}
              step="0.01"
              inputMode="decimal"
              value={Number.isNaN(field.state.value) ? '' : field.state.value}
              placeholder={t('placeholders.price')}
              aria-invalid={field.state.meta.errors.length > 0}
              onBlur={field.handleBlur}
              onChange={(event) =>
                field.handleChange(event.target.valueAsNumber)
              }
            />
            <FieldError errors={toErrors(field.state.meta.errors, 'price')} />
          </Field>
        )}
      />

      <form.Field
        name="category"
        children={(field) => (
          <Field data-invalid={field.state.meta.errors.length > 0}>
            <FieldLabel htmlFor={field.name}>
              {t('fields.category')}
            </FieldLabel>
            <Select
              id={field.name}
              name={field.name}
              value={field.state.value}
              options={categoryOptions}
              aria-invalid={field.state.meta.errors.length > 0}
              onBlur={field.handleBlur}
              onValueChange={field.handleChange}
            />
            <FieldError
              errors={toErrors(field.state.meta.errors, 'category')}
            />
          </Field>
        )}
      />

      <form.Subscribe
        selector={(state) => state.values.category}
        children={(category) =>
          category === CUSTOM_CATEGORY_VALUE && (
            <form.Field
              name="customCategory"
              children={(field) => (
                <Field data-invalid={field.state.meta.errors.length > 0}>
                  <FieldLabel htmlFor={field.name}>
                    {t('fields.customCategory')}
                  </FieldLabel>
                  <Input
                    id={field.name}
                    name={field.name}
                    maxLength={CUSTOM_CATEGORY_MAX}
                    value={field.state.value}
                    placeholder={t('placeholders.customCategory')}
                    aria-invalid={field.state.meta.errors.length > 0}
                    onBlur={field.handleBlur}
                    onChange={(event) =>
                      field.handleChange(event.target.value)
                    }
                  />
                  <FieldError
                    errors={toErrors(field.state.meta.errors, 'customCategory')}
                  />
                </Field>
              )}
            />
          )
        }
      />

      <form.Field
        name="condition"
        children={(field) => (
          <Field data-invalid={field.state.meta.errors.length > 0}>
            <FieldLabel htmlFor={field.name}>
              {t('fields.condition')}
            </FieldLabel>
            <Select
              id={field.name}
              name={field.name}
              value={field.state.value}
              options={conditionOptions}
              aria-invalid={field.state.meta.errors.length > 0}
              onBlur={field.handleBlur}
              onValueChange={field.handleChange}
            />
            <FieldError
              errors={toErrors(field.state.meta.errors, 'condition')}
            />
          </Field>
        )}
      />
    </FieldGroup>
  )
}
