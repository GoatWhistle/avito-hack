import { useEffect, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'

interface LoadMoreProps {
  hasNextPage: boolean
  isFetching: boolean
  onLoadMore: () => void
}

export function LoadMore({
  hasNextPage,
  isFetching,
  onLoadMore,
}: LoadMoreProps) {
  const { t } = useTranslation()
  const sentinelRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const node = sentinelRef.current
    if (!node || !hasNextPage || isFetching) return
    if (typeof IntersectionObserver === 'undefined') return

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries.some((entry) => entry.isIntersecting)) onLoadMore()
      },
      { rootMargin: '200px' },
    )

    observer.observe(node)

    return () => observer.disconnect()
  }, [hasNextPage, isFetching, onLoadMore])

  if (!hasNextPage) return null

  return (
    <div ref={sentinelRef} className="flex justify-center">
      <Button variant="outline" disabled={isFetching} onClick={onLoadMore}>
        {isFetching ? t('status.loading') : t('actions.loadMore')}
      </Button>
    </div>
  )
}
