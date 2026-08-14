import { useEffect, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { useSession } from '#/features/auth/session'
import { FavoriteButton } from '#/features/favorites/components'
import { useFavoriteIds } from '#/features/favorites/hooks'
import {
  useItemPhotosQuery,
  useItemQuery,
  useViewItem,
} from '#/features/items/hooks'
import { formatDate, formatPrice, isFavoritable } from '#/features/items/lib'
import { itemCategories, itemConditions } from '#/features/items/types'
import { ItemSourceBadge } from './ItemSourceBadge'
import { ItemStatusBadge } from './ItemStatusBadge'
import { PhotoGallery } from './PhotoGallery'
import { SoldCelebration } from './SoldCelebration'
import { StatusActions } from './StatusActions'
import { ErrorState, ItemsSkeleton } from './ListStates'

const isKnownCategory = (
  value: string,
): value is (typeof itemCategories)[number] =>
  (itemCategories as readonly string[]).includes(value)

const isKnownCondition = (
  value: string,
): value is (typeof itemConditions)[number] =>
  (itemConditions as readonly string[]).includes(value)

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
  const viewItem = useViewItem()
  const viewedRef = useRef<string | null>(null)

  const loadedItem = itemQuery.data
  const viewerId = user?.id

  useEffect(() => {
    if (!loadedItem) return
    if (viewedRef.current === loadedItem.id) return
    if (viewerId && viewerId === loadedItem.owner_id) return
    viewedRef.current = loadedItem.id
    viewItem.mutate(loadedItem.id)
  }, [loadedItem, viewerId, viewItem])

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
  const rawCategory = item.attributes?.category ?? ''
  const rawCondition = item.attributes?.condition ?? ''
  const categoryLabel = rawCategory
    ? isKnownCategory(rawCategory)
      ? t(`category.${rawCategory}`)
      : rawCategory
    : null
  const conditionLabel =
    rawCondition && isKnownCondition(rawCondition)
      ? t(`condition.${rawCondition}`)
      : null
  const otherAttributes = Object.entries(item.attributes ?? {}).filter(
    ([key]) => key !== 'category' && key !== 'condition',
  )

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
            <ItemSourceBadge
              isSeed={item.is_seed}
              aiVerified={item.ai_verified}
            />
            <span className="text-xs text-muted-foreground">
              {formatDate(item.created_at, i18n.language)}
            </span>
          </div>

          {(categoryLabel ?? conditionLabel) && (
            <p className="text-sm text-muted-foreground">
              {[categoryLabel, conditionLabel].filter(Boolean).join(' · ')}
            </p>
          )}
        </header>

        {celebrating && (
          <SoldCelebration onDismiss={() => setCelebrating(false)} />
        )}

        {isOwner && item.status === 'moderation' && (
          <div
            role="status"
            className="rounded-lg border border-border bg-muted px-3 py-2 text-sm"
          >
            {item.moderation_reason ? (
              <>
                <p className="font-medium text-destructive">
                  {t('moderation.rejectedTitle')}
                </p>
                <p className="text-muted-foreground">
                  {t('moderation.rejectedReason', {
                    reason: item.moderation_reason,
                  })}
                </p>
                <p className="text-muted-foreground">
                  {t('moderation.rejectedHint')}
                </p>
              </>
            ) : (
              <p className="text-muted-foreground">{t('moderation.pending')}</p>
            )}
          </div>
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

        {otherAttributes.length > 0 && (
          <section className="flex flex-col gap-1">
            <h2 className="text-sm font-medium">{t('fields.attributes')}</h2>
            <dl className="grid grid-cols-1 gap-1 text-sm sm:grid-cols-2">
              {otherAttributes.map(([key, value]) => (
                <div key={key} className="flex justify-between gap-2">
                  <dt className="text-muted-foreground">{key}</dt>
                  <dd className="text-right break-words">{value}</dd>
                </div>
              ))}
            </dl>
          </section>
        )}

        {isOwner && item.is_seed && (
          <p className="rounded-lg bg-muted px-3 py-2 text-sm text-muted-foreground">
            {t('seed.locked')}
          </p>
        )}

        {isOwner && !item.is_seed && (
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
