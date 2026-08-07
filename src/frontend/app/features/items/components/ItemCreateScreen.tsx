import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router'
import { Button } from '#/components/ui'
import { useCreateItem } from '#/features/items/hooks'
import { useItemForm } from '#/features/items/forms'
import { ItemFormFields } from './ItemFormFields'
import { QualityHint } from './QualityHint'

export function ItemCreateScreen() {
  const { t } = useTranslation('items')
  const { t: tCommon } = useTranslation()
  const navigate = useNavigate()
  const { mutateAsync } = useCreateItem()
  const [error, setError] = useState<string | null>(null)

  const form = useItemForm({
    onSave: async (payload) => {
      setError(null)
      try {
        const created = await mutateAsync(payload)
        await navigate(`/items/${created.id}/edit`, { replace: true })
      } catch (cause) {
        setError(cause instanceof Error ? cause.message : String(cause))
      }
    },
  })

  return (
    <section className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">{t('actions.create')}</h1>

      <form
        noValidate
        className="flex flex-col gap-6 lg:grid lg:grid-cols-[minmax(0,3fr)_minmax(0,2fr)] lg:items-start"
        onSubmit={(event) => {
          event.preventDefault()
          event.stopPropagation()
          void form.handleSubmit()
        }}
      >
        <div className="flex flex-col gap-6">
          <ItemFormFields form={form} />

          {error && (
            <p role="alert" className="text-sm text-destructive">
              {error}
            </p>
          )}

          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting] as const}
            children={([canSubmit, isSubmitting]) => (
              <div className="flex flex-wrap gap-2">
                <Button type="submit" disabled={!canSubmit || isSubmitting}>
                  {isSubmitting
                    ? tCommon('status.loading')
                    : tCommon('actions.create')}
                </Button>
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => void navigate('/items/mine')}
                >
                  {tCommon('actions.cancel')}
                </Button>
              </div>
            )}
          />

          <p className="text-xs text-muted-foreground">
            {t('photo.saveFirst')}
          </p>
        </div>

        <form.Subscribe
          selector={(state) =>
            [state.values.description, state.values.price] as const
          }
          children={([description, price]) => (
            <QualityHint
              description={description}
              price={price}
              photoCount={0}
            />
          )}
        />
      </form>
    </section>
  )
}
