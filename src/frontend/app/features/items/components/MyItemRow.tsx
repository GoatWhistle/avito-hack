import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { formatDate, formatPrice } from '#/features/items/lib'
import { ItemSourceBadge } from './ItemSourceBadge'
import { ItemStatusBadge } from './ItemStatusBadge'
import { StatusActions } from './StatusActions'
import type { ItemListEntry } from '#/features/items/types'

interface MyItemRowProps {
  item: ItemListEntry
}

export function MyItemRow({ item }: MyItemRowProps) {
  const { t, i18n } = useTranslation('items')
  const { t: tCommon } = useTranslation()

  return (
    <article className="flex flex-col gap-3 rounded-xl bg-card p-4 ring-1 ring-foreground/10">
      <div className="flex flex-wrap items-start justify-between gap-2">
        <div className="flex min-w-0 flex-col gap-1">
          <h3 className="text-sm leading-snug font-medium break-words">
            <Link
              to={`/items/${item.id}`}
              className="rounded-sm outline-none hover:underline focus-visible:ring-3 focus-visible:ring-ring/50"
            >
              {item.title}
            </Link>
          </h3>
          <p className="text-base font-semibold">
            {t('price', { value: formatPrice(item.price, i18n.language) })}
          </p>
          <p className="text-xs text-muted-foreground">
            {formatDate(item.created_at, i18n.language)}
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <ItemStatusBadge status={item.status} />
          <ItemSourceBadge
            isSeed={item.is_seed}
            aiVerified={item.ai_verified}
          />
        </div>
      </div>

      {!item.is_seed && (
        <div className="flex flex-wrap items-center gap-2">
          <Button
            size="sm"
            variant="outline"
            render={<Link to={`/items/${item.id}/edit`} />}
          >
            {tCommon('actions.edit')}
          </Button>
          <StatusActions itemId={item.id} status={item.status} size="sm" />
        </div>
      )}
    </article>
  )
}
