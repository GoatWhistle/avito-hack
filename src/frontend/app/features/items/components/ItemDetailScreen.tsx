import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { useSession } from '#/features/auth/session'
import { FavoriteButton } from '#/features/favorites/components'
import { useFavoriteIds } from '#/features/favorites/hooks'
import { useItemPhotosQuery, useItemQuery } from '#/features/items/hooks'
import { formatDate, formatPrice, isFavoritable } from '#/features/items/lib'
import { ItemStatusBadge } from './ItemStatusBadge'
import { PhotoGallery } from './PhotoGallery'
import { SoldCelebration } from './SoldCelebration'
import { StatusActions } from './StatusActions'
import { ErrorState, ItemsSkeleton } from './ListStates'

interface ItemDetailScreenProps {
  itemId: string
}

export function ItemDetailScreen({ itemId }: ItemDetailScreenProps) {
  const { t, i18n } = useTranslation('items')
  const { t: tCommon } = useTranslation()
  const { user } = useSession()
  const favoriteIds = useFavoriteIds()
  const [celebrating, setCelebrating] = useState(false)
  const [actionError, setActionError] = useState<string | null>(null)

  const itemQuery = useItemQuery(itemId)
  const photosQuery = useItemPhotosQuery(itemId)

  if (itemQuery.isPending) return <ItemsSkeleton count={2} />

  if (itemQuery.isError || !itemQuery.data) {
    return (
      <ErrorState
        message={t('notFound')}
        onRetry={() => void itemQuery.refetch()}
      />
    )
  }

  const item = itemQuery.data
  const isOwner = user?.id === item.owner_id

  return (
    <article className="flex flex-col gap-6 lg:grid lg:grid-cols-[minmax(0,3fr)_minmax(0,2fr)] lg:items-start">
      <PhotoGallery photos={photosQuery.data ?? []} title={item.title} />

      <div className="flex flex-col gap-4">
        <header className="flex flex-col gap-2">
          <div className="flex items-start justify-between gap-3">
            <h1 className="text-xl font-semibold break-words">{item.title}</h1>
            {!isOwner && isFavoritable(item.status) && (
              <FavoriteButton
                itemId={item.id}
                isFavorite={favoriteIds.has(item.id)}
                variant="outline"
                entry={{
                  item_id: item.id,
                  owner_id: item.owner_id,
                  title: item.title,
                  price: item.price,
                  status: item.status,
                }}
              />
            )}
          </div>

          <p className="text-2xl font-semibold">
            {t('price', { value: formatPrice(item.price, i18n.language) })}
          </p>

          <div className="flex flex-wrap items-center gap-2">
            <ItemStatusBadge status={item.status} />
            <span className="text-xs text-muted-foreground">
              {formatDate(item.created_at, i18n.language)}
            </span>
          </div>
        </header>

        {celebrating && (
          <SoldCelebration onDismiss={() => setCelebrating(false)} />
        )}

        {actionError && (
          <p role="alert" className="text-sm text-destructive">
            {actionError}
          </p>
        )}

        {item.description && (
          <section className="flex flex-col gap-1">
            <h2 className="text-sm font-medium">{t('fields.description')}</h2>
            <p className="text-sm leading-relaxed break-words whitespace-pre-line text-muted-foreground">
              {item.description}
            </p>
          </section>
        )}

        {item.attributes && Object.keys(item.attributes).length > 0 && (
          <section className="flex flex-col gap-1">
            <h2 className="text-sm font-medium">{t('fields.attributes')}</h2>
            <dl className="grid grid-cols-1 gap-1 text-sm sm:grid-cols-2">
              {Object.entries(item.attributes).map(([key, value]) => (
                <div key={key} className="flex justify-between gap-2">
                  <dt className="text-muted-foreground">{key}</dt>
                  <dd className="text-right break-words">{value}</dd>
                </div>
              ))}
            </dl>
          </section>
        )}

        {isOwner && (
          <section className="flex flex-col gap-2 border-t border-border pt-4">
            <h2 className="text-sm font-medium">{t('ownerActions')}</h2>
            <div className="flex flex-wrap gap-2">
              <Button
                variant="outline"
                render={<Link to={`/items/${item.id}/edit`} />}
              >
                {tCommon('actions.edit')}
              </Button>
              <StatusActions
                itemId={item.id}
                status={item.status}
                onSold={() => {
                  setActionError(null)
                  setCelebrating(true)
                }}
                onError={setActionError}
              />
            </div>
          </section>
        )}
      </div>
    </article>
  )
}
