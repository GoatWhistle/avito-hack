import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Link, useNavigate } from 'react-router'
import { Button } from '#/components/ui'
import { translateApiError } from '#/api'
import { useSession } from '#/features/auth/session'
import {
  useItemPhotosQuery,
  useItemQuery,
  useUpdateItem,
} from '#/features/items/hooks'
import { useItemForm } from '#/features/items/forms'
import { ItemFormFields } from './ItemFormFields'
import { PhotoManager } from './PhotoManager'
import { QualityHint } from './QualityHint'
import { StatusActions } from './StatusActions'
import { ErrorState, ItemsSkeleton } from './ListStates'

interface ItemEditScreenProps {
  itemId: string
}

export function ItemEditScreen({ itemId }: ItemEditScreenProps) {
  const { t } = useTranslation(['items', 'errors'])
  const { t: tCommon } = useTranslation()
  const navigate = useNavigate()
  const { user } = useSession()
  const [error, setError] = useState<string | null>(null)
  const [saved, setSaved] = useState(false)

  const itemQuery = useItemQuery(itemId)
  const photosQuery = useItemPhotosQuery(itemId)
  const { mutateAsync } = useUpdateItem(itemId)

  const item = itemQuery.data

  const form = useItemForm({
    item,
    onSave: async (payload) => {
      setError(null)
      setSaved(false)
      try {
        await mutateAsync(payload)
        setSaved(true)
      } catch (cause) {
        setError(translateApiError(cause, t) ?? null)
      }
    },
  })

  if (itemQuery.isPending) return <ItemsSkeleton count={2} />

  if (itemQuery.isError || !item) {
    return (
      <ErrorState
        message={t('notFound')}
        onRetry={() => void itemQuery.refetch()}
      />
    )
  }

  if (user && user.id !== item.owner_id) {
    return <ErrorState message={t('notFound')} />
  }

  if (item.is_seed) {
    return <ErrorState message={t('seed.locked')} />
  }

  const photoCount = photosQuery.data?.length ?? 0

  return (
    <section className="flex flex-col gap-6">
      <header className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-xl font-semibold break-words">{item.title}</h1>
        <Button variant="ghost" render={<Link to={`/items/${item.id}`} />}>
          {tCommon('actions.back')}
        </Button>
      </header>

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

          {saved && !error && (
            <p role="status" className="text-sm text-muted-foreground">
              {tCommon('status.saved')}
            </p>
          )}

          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting] as const}
            children={([canSubmit, isSubmitting]) => (
              <div className="flex flex-wrap gap-2">
                <Button type="submit" disabled={!canSubmit || isSubmitting}>
                  {isSubmitting
                    ? tCommon('status.loading')
                    : tCommon('actions.save')}
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

          <PhotoManager itemId={item.id} />
        </div>

        <div className="flex flex-col gap-4">
          <form.Subscribe
            selector={(state) =>
              [state.values.description, state.values.price] as const
            }
            children={([description, price]) => (
              <QualityHint
                description={description}
                price={price}
                photoCount={photoCount}
              />
            )}
          />

          <section className="flex flex-col gap-2 rounded-xl bg-card p-4 ring-1 ring-foreground/10">
            <h2 className="text-sm font-medium">{t('ownerActions')}</h2>
            <StatusActions
              itemId={item.id}
              status={item.status}
              size="sm"
              onError={setError}
            />
          </section>
        </div>
      </form>
    </section>
  )
}
