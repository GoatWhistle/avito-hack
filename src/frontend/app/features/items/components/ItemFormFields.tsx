import { useTranslation } from 'react-i18next'
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
  Input,
} from '#/components/ui'
import { DESCRIPTION_MAX, TITLE_MAX } from '#/features/items/lib'
import { useFieldErrors } from '#/features/items/forms'
import type { useItemForm } from '#/features/items/forms'

interface ItemFormFieldsProps {
  form: ReturnType<typeof useItemForm>
}

export function ItemFormFields({ form }: ItemFormFieldsProps) {
  const { t } = useTranslation('items')
  const toErrors = useFieldErrors()

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
    </FieldGroup>
  )
}
