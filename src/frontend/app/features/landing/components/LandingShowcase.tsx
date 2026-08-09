import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { ItemCard } from '#/features/items/components/ItemCard'
import {
  ErrorState,
  ItemsSkeleton,
} from '#/features/items/components/ListStates'
import { flattenPages, useItemsQuery } from '#/features/items/hooks'
import { SHOWCASE_LIMIT } from '../lib'

export function LandingShowcase() {
  const { t } = useTranslation('landing')
  const query = useItemsQuery({ status: 'published', search: '' })
  const items = flattenPages(query.data?.pages).slice(0, SHOWCASE_LIMIT)

  return (
    <section
      aria-labelledby="landing-showcase"
      className="flex flex-col gap-5 scroll-mt-6"
      id="showcase"
    >
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div className="flex flex-col gap-1">
          <h2
            id="landing-showcase"
            className="font-heading text-xl font-semibold text-foreground sm:text-2xl"
          >
            {t('showcase.title')}
          </h2>
          <p className="max-w-2xl text-sm text-pretty text-muted-foreground">
            {t('showcase.subtitle')}
          </p>
        </div>

        <Button
          variant="outline"
          render={<Link to="/items" />}
          className="no-underline"
        >
          {t('showcase.all')}
        </Button>
      </div>

      {query.isPending && <ItemsSkeleton count={SHOWCASE_LIMIT} />}

      {query.isError && (
        <ErrorState
          message={
            query.error instanceof Error ? query.error.message : undefined
          }
          onRetry={() => void query.refetch()}
        />
      )}

      {query.isSuccess && items.length === 0 && (
        <p
          role="status"
          className="rounded-xl bg-card px-6 py-10 text-center text-sm text-muted-foreground ring-1 ring-foreground/10"
        >
          {t('showcase.empty')}
        </p>
      )}

      {query.isSuccess && items.length > 0 && (
        <ul
          aria-label={t('showcase.title')}
          className="grid list-none auto-rows-fr grid-cols-2 items-stretch gap-3 sm:gap-4 lg:grid-cols-4"
        >
          {items.map((item) => (
            <li key={item.id} className="flex h-full">
              <ItemCard item={item} showFavorite={false} />
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
